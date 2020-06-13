package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/dapr-starter/dapr"
	"github.com/mchmarny/gcputil/env"
	"go.opencensus.io/plugin/ochttp"
)

var (
	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	logger = log.New(os.Stdout, "", 0)

	servicePort    = env.MustGetEnvVar("PORT", "8080")
	subscribeTopic = env.MustGetEnvVar("EVENT_TOPIC_NAME", "events")
	stateStore     = env.MustGetEnvVar("STATE_STORE_NAME", "store")

	daprClient = dapr.NewClient()
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)

	// simple routes
	r.GET("/", rootHandler)

	// pubsub
	r.GET("/dapr/subscribe", subscriptionHandler)
	r.POST("/events", eventHandler)
	r.POST("/message", messagePublisher)

	// server
	hostPort := net.JoinHostPort("0.0.0.0", servicePort)
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := http.ListenAndServe(hostPort, &ochttp.Handler{Handler: r}); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
