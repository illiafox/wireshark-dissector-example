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

![img.png](docs/preview.png)
![img.png](docs/1.png)


#### Client <-> Server Flow

![your-UML-diagram-name](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/illiafox/plantuml-markdown/main/docs/schema.iuml)

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