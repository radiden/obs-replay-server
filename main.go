package main

import (
	"log"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/labstack/echo/v4"
	"github.com/radiden/obs-replay-server/handlers"
)

func main() {
	// init obs connection
	client, err := goobs.New("localhost:4455", goobs.WithPassword("ThpNdGXBwneYn9X7"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	isRecording, err := client.Outputs.GetReplayBufferStatus()
	if err != nil {
		panic(err)
	}
	if !isRecording.OutputActive {
		client.Outputs.StartReplayBuffer()
	}

	// init echo
	e := echo.New()
	e.GET("/", handlers.IndexHandler)
	e.Logger.Fatal(e.Start(":1323"))
}

func getSavedPath(client *goobs.Client) string {
	for {
		msg := <-client.IncomingEvents
		switch m := msg.(type) {
		case *events.ReplayBufferSaved:
			return m.SavedReplayPath
		}
	}
}
