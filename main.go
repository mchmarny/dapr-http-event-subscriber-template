package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	// Version will be set during build
	Version = "v0.0.1-default"

	logger = log.New(os.Stdout, "", 0)

	servicePort = getEnvVar("PORT", "8080")
	topicName   = getEnvVar("TOPIC_NAME", "events")
	storeName   = getEnvVar("STORE_NAME", "store")

	// dapr
	daprClient dapr.Client
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// wire actual Dapr client
	c, err := dapr.NewClient()
	if err != nil {
		logger.Fatalf("error creating Dapr client: %v", err)
	}
	daprClient = c

	// router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)

	// pubsub
	r.GET("/dapr/subscribe", subscriptionHandler)
	r.POST("/events", eventHandler)

	// default route
	r.Any("/", defaultHandler)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", Version, hostPort)
	if err := http.ListenAndServe(hostPort, r); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}

// Options midleware
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "POST,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

func getEnvVar(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(val)
	}
	return fallbackValue
}
