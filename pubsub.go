package main

import (
	"encoding/json"
	"net/http"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
)

const (
	// SupportedCloudEventVersion indicates the version of CloudEvents suppored by this handler
	SupportedCloudEventVersion = "0.3"

	//SupportedCloudEventContentTye indicates the content type supported by this handlers
	SupportedCloudEventContentTye = "application/json"
)

// DaprSubscription represents single Dapr subscription
type DaprSubscription struct {
	Topic string `json:"topic"`
	Route string `json:"route"`
}

func subscriptionHandler(c *gin.Context) {
	subscriptions := []DaprSubscription{
		{
			Topic: subscribeTopic,
			Route: "/events",
		},
	}
	logger.Printf("subscription topics: %v", subscriptions)
	c.JSON(http.StatusOK, subscriptions)
}

func eventHandler(c *gin.Context) {
	ctx := getTraceContext(c)
	e := ce.NewEvent()
	if err := c.ShouldBindJSON(&e); err != nil {
		logger.Printf("error binding event: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	eventVersion := e.Context.GetSpecVersion()
	if eventVersion != SupportedCloudEventVersion {
		logger.Printf("invalid event spec version: %s", eventVersion)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	eventContentType := e.Context.GetDataContentType()
	if eventContentType != SupportedCloudEventContentTye {
		logger.Printf("invalid event content type: %s", eventContentType)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	msgIn := SimpleMessage{}
	err := json.Unmarshal(e.Data(), &msgIn)
	if err != nil {
		logger.Printf("error parsing message from event: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	logger.Printf("saving event %s to %s - %v", e.ID(), stateStore, msgIn)
	err = daprClient.SaveState(ctx, stateStore, e.ID(), msgIn)
	if err != nil {
		logger.Printf("error saving event to store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	out, err := daprClient.GetState(ctx, stateStore, e.ID())
	if err != nil {
		logger.Printf("error retreaving event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	msgOut := SimpleMessage{}
	err = json.Unmarshal(out, &msgOut)
	if err != nil {
		logger.Printf("error parsing event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("retreaved event %s from %s - %v", e.ID(), stateStore, msgOut)

	err = daprClient.DeleteState(ctx, stateStore, e.ID())
	if err != nil {
		logger.Printf("error deleting event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("event %s deleted from %s", e.ID(), stateStore)

	c.JSON(http.StatusOK, nil)
}

// SimpleMessage is a message posted to publsher
type SimpleMessage struct {
	Message string `json:"message"`
}

func messagePublisher(c *gin.Context) {
	ctx := getTraceContext(c)
	m := SimpleMessage{}
	if err := c.ShouldBindJSON(&m); err != nil {
		logger.Printf("error binding message: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	b, err := json.Marshal(m)
	if err != nil {
		logger.Printf("error serializing message: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	err = daprClient.Publish(ctx, subscribeTopic, b)
	if err != nil {
		logger.Printf("error publishing event: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	c.JSON(http.StatusOK, nil)
}
