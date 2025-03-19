package proto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// | OpCode | ID Length | Station ID  | Temperature | Humidity | Pressure |
// | ------ | --------- | ----------- | ----------- | -------- | -------- |
// | 1 byte | 1 byte    | (id length) | 4 bytes     | 4 bytes  | 4 bytes  |

// SendMessage represents data sent from weather station
// Endianness: Little Endian
//
// 1 byte - ID length in bytes
// [ID length] - ID
// 4 bytes - Temperature
// 4 bytes - Humidity
// 4 bytes - Pressure
type SendMessage struct {
	ID          string
	Temperature float32
	Humidity    float32
	Pressure    float32
}

func (SendMessage) OpCode() OpCode { return SEND }

func (s SendMessage) Binary() ([]byte, error) {
	buf := new(bytes.Buffer)

	idLength := len(s.ID)
	if idLength > 255 {
		return nil, errors.New("id is too long")
	}

	err := binary.Write(buf, binary.LittleEndian, byte(idLength))
	if err != nil {
		return nil, fmt.Errorf("write id length %d: %w", idLength, err)
	}

	buf.WriteString(s.ID)

	err = binary.Write(buf, binary.LittleEndian, s.Temperature)
	if err != nil {
		return nil, fmt.Errorf("write temperature %f: %w", s.Temperature, err)
	}

	err = binary.Write(buf, binary.LittleEndian, s.Humidity)
	if err != nil {
		return nil, fmt.Errorf("write humidity %f: %w", s.Humidity, err)
	}

	err = binary.Write(buf, binary.LittleEndian, s.Pressure)
	if err != nil {
		return nil, fmt.Errorf("write pressure %f: %w", s.Pressure, err)
	}

	return buf.Bytes(), nil
}

func (s *SendMessage) ReadBinary(reader io.Reader) error {

	// Read ID length (1 byte)
	var idLength byte
	err := binary.Read(reader, binary.LittleEndian, &idLength)
	if err != nil {
		return fmt.Errorf("failed to read id length: %w", err)
	}

	// Read the ID (ID length bytes)
	id := make([]byte, idLength)
	err = binary.Read(reader, binary.LittleEndian, &id)
	if err != nil {
		return fmt.Errorf("failed to read ID: %w", err)
	}

	// Read Temperature (4 bytes)
	var temperature float32
	err = binary.Read(reader, binary.LittleEndian, &temperature)
	if err != nil {
		return fmt.Errorf("failed to read temperature: %w", err)
	}

	// Read Humidity (4 bytes)
	var humidity float32
	err = binary.Read(reader, binary.LittleEndian, &humidity)
	if err != nil {
		return fmt.Errorf("failed to read humidity: %w", err)
	}

	// Read Pressure (4 bytes)
	var pressure float32
	err = binary.Read(reader, binary.LittleEndian, &pressure)
	if err != nil {
		return fmt.Errorf("failed to read pressure: %w", err)
	}

	s.ID = string(id)
	s.Temperature = temperature
	s.Humidity = humidity
	s.Pressure = pressure

	return nil
}
