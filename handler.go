package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

	logger.Printf("saving event %s to %s - %s", e.ID(), storeName, string(e.Data()))
	err := daprClient.SaveState(ctx, storeName, e.ID(), e.Data())
	if err != nil {
		logger.Printf("error saving event to store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}

	out, err := daprClient.GetState(ctx, storeName, e.ID())
	if err != nil {
		logger.Printf("error retreaving event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("retreaved event %s from %s - %s", e.ID(), storeName, string(out))

	err = daprClient.DeleteState(ctx, storeName, e.ID())
	if err != nil {
		logger.Printf("error deleting event from store: %v", err)
		c.JSON(http.StatusBadRequest, clientError)
		return
	}
	logger.Printf("event %s deleted from %s", e.ID(), storeName)

	c.JSON(http.StatusOK, nil)
}
