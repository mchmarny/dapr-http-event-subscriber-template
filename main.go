package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"

	dapr "github.com/mchmarny/godapr/v1"
)

var (
	// AppVersion will be overritten during build
	AppVersion = "v0.0.1-default"

	logger = log.New(os.Stdout, "", 0)

	servicePort = env.MustGetEnvVar("PORT", "8080")
	topicName   = env.MustGetEnvVar("TOPIC_NAME", "events")
	storeName   = env.MustGetEnvVar("STORE_NAME", "store")

	// dapr
	daprClient Client

	// test client against local interace
	_ = Client(dapr.NewClient())
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// wire actual Dapr client
	daprClient = dapr.NewClient()

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
	logger.Printf("Server (%s) starting: %s \n", AppVersion, hostPort)
	if err := http.ListenAndServe(hostPort, &ochttp.Handler{Handler: r}); err != nil {
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

func getTraceContext(c *gin.Context) trace.SpanContext {
	httpFmt := tracecontext.HTTPFormat{}
	ctx, ok := httpFmt.SpanContextFromRequest(c.Request)
	if !ok {
		ctx = trace.SpanContext{}
	}

	logger.Printf("trace info [%s]: 0-%x-%x-%x",
		c.Request.URL.Path,
		ctx.TraceID[:],
		ctx.SpanID[:],
		[]byte{byte(ctx.TraceOptions)})

	return ctx
}

// Client is the minimal client support for testing
type Client interface {
	GetState(ctx trace.SpanContext, store, key string) ([]byte, error)
	SaveState(ctx trace.SpanContext, store, key string, data interface{}) error
	DeleteState(ctx trace.SpanContext, store, key string) error
}
