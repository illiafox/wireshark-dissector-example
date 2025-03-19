package proto

// AckMessage represents an acknowledgment from server after receiving data
type AckMessage struct{}

func (AckMessage) OpCode() OpCode { return ACK }

func (s AckMessage) Binary() ([]byte, error) {
	return nil, nil
}
