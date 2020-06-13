package dapr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientCreation(t *testing.T) {
	url := "http://localhost:3500"
	c := NewClient()
	assert.NotNil(t, c)
	assert.Equal(t, c.url, url)

	c1 := NewClientWithURL(url)
	assert.NotNil(t, c1)
	assert.Equal(t, c.url, url)
}

// func TestState(t *testing.T) {
// 	c := NewClient()
// 	assert.NotNil(t, c)

// 	ctx := context.Background()
// 	data := &Person{
// 		Name: "Test Person",
// 		Age:  40,
// 	}

// 	err := c.SaveState(ctx, "statestore", "key1", data)
// 	assert.Nil(t, err)

// 	out, err := c.GetState(ctx, "statestore", "key1")
// 	assert.Nil(t, err)
// 	assert.NotNil(t, out)

// 	data2 := Person{}
// 	err = json.Unmarshal(out, &data2)
// 	assert.NotNil(t, data2)
// 	assert.Equal(t, data.Name, data2.Name)

// 	err = c.DeleteState(ctx, "statestore", "key1")
// 	assert.Nil(t, err)
// }

// // Person is a test object for this example
// type Person struct {
// 	Name string
// 	Age  int
// }

// // EchoMessage holds the request and response
// type EchoMessage struct {
// 	Request  string `json:"req"`
// 	Response string `json:"res"`
// }
