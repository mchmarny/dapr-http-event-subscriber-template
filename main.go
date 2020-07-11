package main

import (
	"context"

	"log"
	"net/http"
	"os"
	"strings"

	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	logger      = log.New(os.Stdout, "", 0)
	servicePort = getEnvVar("PORT", "8080")
	topicName   = getEnvVar("TOPIC_NAME", "events")
)

func main() {
	// create a regular HTTP server mux
	mux := http.NewServeMux()

	// create a Dapr service
	s, err := daprd.NewService(mux)
	if err != nil {
		logger.Fatalf("error creating sever: %v", err)
	}

	// add some topic subscriptions
	err = s.AddTopicEventHandler("events", "/events", eventHandler)
	if err != nil {
		logger.Fatalf("error adding topic subscription: %v", err)
	}

	// start the server
	err = s.HandleSubscriptions()
	if err != nil {
		logger.Fatalf("error creating topic subscription: %v", err)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error listenning: %v", err)
	}
}

func eventHandler(ctx context.Context, e daprd.TopicEvent) error {
	logger.Printf("event - Source: %s, Topic:%s, ID:%s", e.Source, e.Topic, e.ID)

	switch v := e.Data.(type) {
	case string:
		logger.Printf("%s", e.Data.(string))
	case map[string]interface{}:
		c := e.Data.(map[string]interface{})
		for k, v := range c {
			logger.Printf("%s: %v", k, v)
		}
	default:
		logger.Printf("%t", v)
	}

	// TODO: do something with the cloud event data

	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
