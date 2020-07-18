package main

import (
	"context"
	"fmt"
	"time"

	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	logger    *log.Logger
	address   = getEnvVar("ADDRESS", ":8080")
	topicName = getEnvVar("TOPIC_NAME", "events")
)

func init() {
	// configure logging
	logger = log.New()
	logger.Level = log.DebugLevel
	logger.Out = os.Stdout
	logger.Formatter = &log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyLevel: "severity",
			log.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
}

func main() {
	// create a Dapr service
	s := daprd.NewService()

	// add some topic subscriptions
	topicRoute := fmt.Sprintf("/%s", topicName)
	err := s.AddTopicEventHandler(topicName, topicRoute, eventHandler)
	if err != nil {
		logger.Fatalf("error adding topic subscription: %v", err)
	}

	// start the service
	if err = s.Start(address); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("error starting service: %v", err)
	}
}

func eventHandler(ctx context.Context, e daprd.TopicEvent) error {
	logger.Debugf(
		"event - Source: %s, Topic:%s, ID:%s, DataContentType:%s",
		e.Source, e.Topic, e.ID, e.DataContentType,
	)

	// TODO: do something with the cloud event data
	logger.Infoln(e.Data)

	return nil
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
