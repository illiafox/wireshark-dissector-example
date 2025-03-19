package proto

import "bytes"

type Message interface {
	Binary() ([]byte, error)
	OpCode() OpCode
}

func Marshall(m Message) ([]byte, error) {
	data, err := m.Binary()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Grow(1 + len(data)) // + 1 for opcode
	buf.WriteByte(byte(m.OpCode()))
	buf.Write(data)

	return buf.Bytes(), nil
}
