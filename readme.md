# Wireshark Dissector Example

1. Start server
```bash
go run server/server.go
```

2. Start client
```bash
go run ./client --addr=localhost:12345
```

3. Copy [weather.lua](./wireshark/dissectors/weather.lua) to `%WIRESHARK_PATH%/plugins/weather/`

#### Screenshots

![img.png](screenshots/preview.png)
![img.png](screenshots/1.png)


#### Client <-> Server Flow

```plantuml
@startuml
Client -> Server : SEND (0x2): ID, Temperature, Humidity, Pressure
Server -> Client : ACK (0x1)
Client -> Server : SEND (0x2): ID, Temperature, Humidity, Pressure
Server -> Client : ACK (0x1)
Client -> Server : SEND (0x2): ID, Temperature, Humidity, Pressure
Server -> Client : ACK (0x1)
@enduml
```