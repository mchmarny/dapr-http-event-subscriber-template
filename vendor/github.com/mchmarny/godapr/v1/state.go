package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// SaveStateWithData saves state data into state store
func (c *Client) SaveStateWithData(ctx trace.SpanContext, store string, data *StateData) error {
	if store == "" {
		return errors.New("nil store")
	}
	if data == nil {
		return errors.New("nil input data")
	}

	list := []*StateData{data}
	url := fmt.Sprintf("%s/v1.0/state/%s", c.url, store)
	b, _ := json.Marshal(list)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	_, status, err := c.exec(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "error quering state service: %s", url)
	}

	if status != http.StatusCreated {
		return fmt.Errorf("invalid response code to %s: %d", url, status)
	}

	return nil
}

// SaveState saves data into state store for specific key
func (c *Client) SaveState(ctx trace.SpanContext, store, key string, data interface{}) error {
	state := &StateData{
		Key:     key,
		Value:   data,
		Options: DefaultStateOptions,
		Metadata: map[string]string{
			"created_on": time.Now().UTC().String(),
		},
	}
	return c.SaveStateWithData(ctx, store, state)
}

// GetStateWithOptions gets content for specific key in state store
func (c *Client) GetStateWithOptions(ctx trace.SpanContext, store, key string, opt *StateOptions) (data []byte, err error) {
	if opt == nil {
		return nil, errors.New("nil state options")
	}

	url := fmt.Sprintf("%s/v1.0/state/%s/%s", c.url, store, key)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if opt.Concurrency != "" {
		req.Header.Set("concurrency", opt.Concurrency)
	}
	if opt.Consistency != "" {
		req.Header.Set("consistency", opt.Consistency)
	}

	content, status, err := c.exec(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "error quering state service: %s", url)
	}

	// on initial run there won't be any state
	if status == http.StatusNoContent || status == http.StatusNotFound {
		return nil, nil
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("invalid response code from GET to %s: %d", url, status)
	}

	return content, nil
}

// GetState gets content for specific key in state store
func (c *Client) GetState(ctx trace.SpanContext, store, key string) (data []byte, err error) {
	return c.GetStateWithOptions(ctx, store, key, DefaultStateOptions)
}

// DeleteStateWithOptions deletes existing state from specified store
func (c *Client) DeleteStateWithOptions(ctx trace.SpanContext, store, key string, opt *StateOptions) error {
	if opt == nil {
		return errors.New("nil state options")
	}

	url := fmt.Sprintf("%s/v1.0/state/%s/%s", c.url, store, key)
	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if opt.Concurrency != "" {
		req.Header.Set("concurrency", opt.Concurrency)
	}

	if opt.Consistency != "" {
		req.Header.Set("consistency", opt.Consistency)
	}

	_, status, err := c.exec(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "error quering state service: %s", url)
	}

	// on initial run there won't be any state
	if status == http.StatusNoContent || status == http.StatusNotFound {
		return nil
	}

	if status != http.StatusOK {
		return fmt.Errorf("invalid response code from GET to %s: %d", url, status)
	}

	return nil
}

// DeleteState deletes existing state from specified store
func (c *Client) DeleteState(ctx trace.SpanContext, store, key string) error {
	return c.DeleteStateWithOptions(ctx, store, key, DefaultStateOptions)
}
