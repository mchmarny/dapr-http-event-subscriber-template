package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	dapr "github.com/mchmarny/godapr/v1"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"
)

func TestDefaultHandler(t *testing.T) {
	daprClient = GetTestClient()
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)
	r.GET("/", defaultHandler)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSubscriptionHandler(t *testing.T) {
	daprClient = GetTestClient()
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)
	r.GET("/", subscriptionHandler)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	content, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)

	var subs []dapr.Subscription
	err = json.Unmarshal(content, &subs)
	assert.Nil(t, err)
	assert.Lenf(t, subs, 1, "minimum 1 subscription required, got: %v", subs)
}

func TestEventHandler(t *testing.T) {
	daprClient = GetTestClient()
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)
	r.POST("/", eventHandler)
	w := httptest.NewRecorder()

	data, err := ioutil.ReadFile("./event.json")
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func GetTestClient() *TestClient {
	return &TestClient{}
}

var (
	// test test client against local interace
	_ = Client(&TestClient{})
)

type TestClient struct {
}

func (c *TestClient) GetState(ctx trace.SpanContext, store, key string) ([]byte, error) {
	return []byte(`{ "message": "hello" }`), nil
}
func (c *TestClient) SaveState(ctx trace.SpanContext, store, key string, data interface{}) error {
	return nil
}
func (c *TestClient) DeleteState(ctx trace.SpanContext, store, key string) error {
	return nil
}
