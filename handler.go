package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	dapr "github.com/mchmarny/godapr/v1"

	ce "github.com/cloudevents/sdk-go/v2"
)

const (
	// SupportedCloudEventVersion indicates the version of CloudEvents suppored by this handler
	SupportedCloudEventVersion = "0.3"

	//SupportedCloudEventContentTye indicates the content type supported by this handlers
	SupportedCloudEventContentTye = "application/json"
)

var (
	clientError = gin.H{
		"error":   "Bad Request",
		"message": "Error processing your request, see logs for details",
	}
)

// SimpleMessage corresponds to the payload published in make event
type SimpleMessage struct {
	Message string `json:"message"`
}

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      AppVersion,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func subscriptionHandler(c *gin.Context) {
	subscriptions := []dapr.Subscription{
		{
			Topic: topicName,
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

	var in SimpleMessage
	if err := json.Unmarshal(e.Data(), &in); err != nil {
		logger.Printf("invalid event content format in event: %s", string(e.Data()))
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	logger.Printf("saving event %s to %s - %v", e.ID(), storeName, in)
	err := daprClient.SaveState(ctx, storeName, e.ID(), in)
	if err != nil {
		logger.Printf("error saving event to store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	state, err := daprClient.GetState(ctx, storeName, e.ID())
	if err != nil {
		logger.Printf("error retreaving event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	var out SimpleMessage
	if err := json.Unmarshal(state, &out); err != nil {
		logger.Printf("invalid event content format from store: %s", string(state))
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("retreaved event %s from %s - %v", e.ID(), storeName, out)

	err = daprClient.DeleteState(ctx, storeName, e.ID())
	if err != nil {
		logger.Printf("error deleting event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("event %s deleted from %s", e.ID(), storeName)

	c.JSON(http.StatusOK, nil)
}
