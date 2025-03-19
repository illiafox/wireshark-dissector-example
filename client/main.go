package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/illiafox/wireshark-dissector-example/proto"
)

func main() {
	addr := flag.String("addr", "localhost:6080", "server address")
	interval := flag.Duration("interval", 1*time.Second, "send interval")
	stationName := flag.String("station", "meteo-test", "station name")

	flag.Parse()

	client, err := NewClient(*addr)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	for range time.Tick(*interval) {
		temperature := rand.Float32() * 25
		humidity := rand.Float32() * 80
		pressure := rand.Float32() * 100

		msg := proto.SendMessage{
			ID:          *stationName,
			Temperature: temperature,
			Humidity:    humidity,
			Pressure:    pressure,
		}

		fmt.Printf("Sending message: %+v\n", msg)

		err = client.SendMessage(ctx, msg)
		if err != nil {
			log.Println("failed to send message", err)
		}
	}
}
