package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.ReleaseMode)
	daprClient = &TestClient{}
	r := m.Run()
	os.Exit(r)
}

func TestDefaultHandler(t *testing.T) {
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
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)
	r.GET("/", subscriptionHandler)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	_, err := ioutil.ReadAll(w.Body)
	assert.Nil(t, err)
}

func TestEventHandler(t *testing.T) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Options)
	r.POST("/", eventHandler)
	w := httptest.NewRecorder()

	data, err := ioutil.ReadFile("./event/sample.json")
	assert.Nil(t, err)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

var (
	// test test client against local interace
	_ = dapr.Client(&TestClient{})
)

type TestClient struct {
}

func (c *TestClient) InvokeBinding(ctx context.Context, name, op string, in []byte, min map[string]string) (out []byte, mout map[string]string, err error) {
	return []byte("hello"), make(map[string]string), nil
}

func (c *TestClient) InvokeOutputBinding(ctx context.Context, name, operation string, data []byte) error {
	return nil
}

func (c *TestClient) InvokeService(ctx context.Context, serviceID, method string) (out []byte, err error) {
	return []byte("hello"), nil
}

func (c *TestClient) InvokeServiceWithContent(ctx context.Context, serviceID, method, contentType string, data []byte) (out []byte, err error) {
	return []byte("hello"), nil
}

func (c *TestClient) PublishEvent(ctx context.Context, topic string, in []byte) error {
	return nil
}

func (c *TestClient) GetSecret(ctx context.Context, store, key string, meta map[string]string) (out map[string]string, err error) {
	return make(map[string]string), nil
}

func (c *TestClient) SaveState(ctx context.Context, s *dapr.State) error {
	return nil
}

func (c *TestClient) SaveStateData(ctx context.Context, store, key string, data []byte) error {
	return nil
}

func (c *TestClient) SaveStateDataVersion(ctx context.Context, store, key, etag string, data []byte) error {
	return nil
}

func (c *TestClient) SaveStateItem(ctx context.Context, store string, item *dapr.StateItem) error {
	return nil
}

func (c *TestClient) GetState(ctx context.Context, store, key string) (out []byte, etag string, err error) {
	return []byte("hello"), "", nil
}

func (c *TestClient) GetStateWithConsistency(ctx context.Context, store, key string, sc dapr.StateConsistency) (out []byte, etag string, err error) {
	return []byte("hello"), "", nil
}

func (c *TestClient) DeleteState(ctx context.Context, store, key string) error {
	return nil
}

func (c *TestClient) DeleteStateVersion(ctx context.Context, store, key, etag string, opts *dapr.StateOptions) error {
	return nil
}

func (c *TestClient) Close() {

}
