package proto

import (
	"bytes"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendMessage_Binary(t *testing.T) {
	temperature := rand.Float32() * 25
	humidity := rand.Float32() * 80
	pressure := rand.Float32() * 100

	name := "test1234"

	msg := SendMessage{
		ID:          name,
		Temperature: temperature,
		Humidity:    humidity,
		Pressure:    pressure,
	}

	data, err := msg.Binary()
	require.NoError(t, err, "marshall binary")

	var parse SendMessage
	err = parse.ReadBinary(bytes.NewReader(data))
	require.NoError(t, err, "unmarshall binary")

	require.Equal(t, msg, parse)
}
