package main

import (
	"context"
	"encoding/json"

	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dapr/go-sdk/server/event"
	daprd "github.com/dapr/go-sdk/server/http"
	"github.com/pkg/errors"
)

var (
	logger      = log.New(os.Stdout, "", 0)
	servicePort = getEnvVar("PORT", "8080")
	topicName   = getEnvVar("TOPIC_NAME", "events")
)

func main() {
	// create a regular HTTP server mux
	mux := http.NewServeMux()

	// create a Dapr service server
	daprServer, err := daprd.NewServer(mux)
	if err != nil {
		log.Fatalf("error creating sever: %v", err)
	}

	// add some topic subscriptions
	err = daprServer.AddTopicEventHandler("events", "/events", eventHandler)
	if err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	// start the server
	err = daprServer.HandleSubscriptions()
	if err != nil {
		log.Fatalf("error creating topic subscription: %v", err)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error listenning: %v", err)
	}
}

func eventHandler(ctx context.Context, e event.TopicEvent) error {
	log.Printf("event - Topic:%s, ID:%s", e.Topic, e.ID)

	var cloudevent CloudEvent
	if err := json.Unmarshal(e.Data, &cloudevent); err != nil {
		return errors.Wrap(err, "error binding cloud event")
	}
	logger.Printf("event: %v", cloudevent.Data)

	// TODO: do something with the cloud event data

	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}

// CloudEvent is a local copy of the minimal Cloud Event message
type CloudEvent struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	SpecVersion     string      `json:"specversion"`
	DataContentType string      `json:"datacontenttype"`
	Data            interface{} `json:"data"`
	Subject         string      `json:"subject"`
}
