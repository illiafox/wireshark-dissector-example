weather_protocol = Proto("weather_protocol", "Weather Protocol")

local f_opcode = ProtoField.uint8("weather_protocol.opcode", "OpCode", base.HEX)
local f_id_length = ProtoField.uint8("weather_protocol.id_length", "ID Length", base.DEC)
local f_id = ProtoField.string("weather_protocol.id", "ID")
local f_temperature = ProtoField.float("weather_protocol.temperature", "Temperature", base.DEC)
local f_humidity = ProtoField.float("weather_protocol.humidity", "Humidity", base.DEC)
local f_pressure = ProtoField.float("weather_protocol.pressure", "Pressure", base.DEC)

weather_protocol.fields = { f_opcode, f_id_length, f_id, f_temperature, f_humidity, f_pressure }

-- OpCodes
local ACK = 0x01
local SEND = 0x02

-- Dissector function
function weather_protocol.dissector(buffer, pinfo, tree)
    local subtree = tree:add(weather_protocol, buffer(), "Weather Protocol Data")

    local opCode = buffer(0, 1):uint()

    pinfo.cols.protocol = "Weather Protocol"

    if opCode == SEND then
        -- 1 byte: ID length
        -- [ID length] bytes: ID
        -- 4 bytes: Temperature
        -- 4 bytes: Humidity
        -- 4 bytes: Pressure

        local id_length = buffer(1, 1):uint()
        subtree:add(f_opcode, opCode)
        subtree:add(f_id_length, id_length)

        local id = buffer(2, id_length):string()
        subtree:add(f_id, id)

        temperature = buffer(2 + id_length, 4):le_float()
        humidity = buffer(6 + id_length, 4):le_float()
        pressure = buffer(10 + id_length, 4):le_float()

        subtree:add(f_temperature, temperature)
        subtree:add(f_humidity, humidity)
        subtree:add(f_pressure, pressure)

        pinfo.cols.info = string.format("--> SEND: Station '%s', Temperature %f, Humidity %f, Pressure %f",
                id, temperature, humidity, pressure
        )

    elseif opCode == ACK then
        pinfo.cols.info = "<-- ACK from server"

        -- AckMessage format:
        -- 1 byte: OpCode
        subtree:add(f_opcode, buffer(0, 1))
    else
        -- Unknown or unhandled OpCode
        subtree:add(f_opcode, buffer(0, 1), "Unknown OpCode: " .. opCode)
    end
end

tcp_table = DissectorTable.get("tcp.port")
tcp_table:add(12345, weather_protocol)