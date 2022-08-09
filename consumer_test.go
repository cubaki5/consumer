package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func makeRequest(batch Batch) ([]byte, error) {
	b, err := json.Marshal(batch)
	if err != nil {
		return nil, err
	}

	rsp, err := http.Post("http://localhost:1323/", "json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	bs, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func TestConsumer(t *testing.T) {
	t.Run("When sending 3 Items", func(t *testing.T) {
		var batch Batch = make([]Item, 3)
		actualRequestBody, err := makeRequest(batch)
		assert.NoError(t, err)
		assert.Equal(t, "Ok", string(actualRequestBody))
	})
	t.Run("Server is full", func(t *testing.T) {
		t.Run("When sending 6 Items", func(t *testing.T) {
			var batch Batch = make([]Item, 6)
			actualRequestBody, err := makeRequest(batch)
			assert.NoError(t, err)
			assert.Equal(t, "Server is full", string(actualRequestBody))
		})
		t.Run("When sending Items after 6 Items without waiting", func(t *testing.T) {
			var batch Batch = make([]Item, 6)
			_, err := makeRequest(batch)
			assert.NoError(t, err)
			actualRequestBody, err := makeRequest(batch)
			assert.NoError(t, err)
			assert.Equal(t, "Server is full", string(actualRequestBody))
		})
		t.Run("When sending 3 Items after 6 Items with 10 seconds waiting", func(t *testing.T) {
			var firstBatch Batch = make([]Item, 6)
			_, err := makeRequest(firstBatch)
			assert.NoError(t, err)
			time.Sleep(10 * time.Second)
			var secondBatch Batch = make([]Item, 3)
			actualRequestBody, err := makeRequest(secondBatch)
			assert.NoError(t, err)
			assert.Equal(t, "Ok", string(actualRequestBody))
		})
		t.Run("When sending 6 Items after 6 Items with 10 seconds waiting", func(t *testing.T) {
			var batch Batch = make([]Item, 6)
			_, err := makeRequest(batch)
			assert.NoError(t, err)
			time.Sleep(10 * time.Second)
			actualRequestBody, err := makeRequest(batch)
			assert.NoError(t, err)
			assert.Equal(t, "Server is full", string(actualRequestBody))
		})
	})
}

func TestConsumer_logs(t *testing.T) {
	t.Log("Given the nees to test Consumer  at different time.")
	{
		testID := 0
		t.Logf("\tTest %d:\t When sending 3 Items", testID)
		{
			var batch Batch = make([]Item, 3)
			actualRequestBody, err := makeRequest(batch)
			assert.NoError(t, err)
			assert.Equal(t, "Ok", string(actualRequestBody))
		}
	}
}
