package main

import (
	"net/http"
	"time"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
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

func defaultHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"release":      Version,
		"request_on":   time.Now(),
		"request_from": c.Request.RemoteAddr,
	})
}

func subscriptionHandler(c *gin.Context) {
	//TODO: use dapr.Subscription when https://github.com/dapr/go-sdk/pull/27 lands
	subscriptions := []*subscription{
		{
			Topic: topicName,
			Route: "/events",
		},
	}
	logger.Printf("subscription topics: %v", subscriptions)
	c.JSON(http.StatusOK, subscriptions)
}

func eventHandler(c *gin.Context) {
	e := ce.NewEvent()
	if err := c.ShouldBindJSON(&e); err != nil {
		logger.Printf("error binding event: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("event: %v", e)

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

	logger.Printf("saving event %s to %s", e.ID(), storeName)
	err := daprClient.SaveStateData(c.Request.Context(), storeName, e.ID(), e.Data())
	if err != nil {
		logger.Printf("error saving event to store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	logger.Printf("geting event %s from %s", e.ID(), storeName)
	out, etag, err := daprClient.GetState(c.Request.Context(), storeName, e.ID())
	if err != nil {
		logger.Printf("error retreaving event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("retreaved event (etag: %s) from %s - %s", etag, storeName, string(out))

	logger.Printf("deleting event %s from %s", e.ID(), storeName)
	err = daprClient.DeleteState(c.Request.Context(), storeName, e.ID())
	if err != nil {
		logger.Printf("error deleting event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Println("event processing done")
	c.JSON(http.StatusOK, nil)
}

type subscription struct {
	Topic string `json:"topic"`
	Route string `json:"route"`
}
