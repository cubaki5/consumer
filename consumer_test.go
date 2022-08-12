package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"consumer/models"
	"consumer/module/consumer"
)

func makeRequest(batch models.Batch) ([]byte, error) {
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
	t.Run("Service tests", func(t *testing.T) {
		tests := []struct {
			name   string
			batch  models.Batch
			expMsg string
		}{
			{
				name:   "When sending batch of len that less then buffer",
				batch:  make([]models.Item, 3),
				expMsg: "Ok",
			},
			{
				name:   "When sending batch of len that equal buffer",
				batch:  make([]models.Item, 5),
				expMsg: "Ok",
			},
			{
				name:   "When sending batch of len that more than buffer",
				batch:  make([]models.Item, 6),
				expMsg: "server is full",
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				actualRequestBody, err := makeRequest(test.batch)
				require.NoError(t, err)
				assert.Equal(t, test.expMsg, actualRequestBody)
			})
		}
		t.Run("When sending Batch after panic without waiting", func(t *testing.T) {
			var batch models.Batch = make([]models.Item, 6)
			_, err := makeRequest(batch)
			require.NoError(t, err)
			actualRequestBody, err := makeRequest(batch)
			require.NoError(t, err)
			assert.Equal(t, "server is full", string(actualRequestBody))
		})
		t.Run("When sending Batch of len that equal buffer after panic with waiting set off panic", func(t *testing.T) {
			var firstBatch models.Batch = make([]models.Item, 6)
			_, err := makeRequest(firstBatch)
			require.NoError(t, err)
			time.Sleep(consumer.PanicDuration)
			var secondBatch models.Batch = make([]models.Item, 5)
			actualRequestBody, err := makeRequest(secondBatch)
			require.NoError(t, err)
			require.Equal(t, "Ok", string(actualRequestBody))
		})
	})
}
