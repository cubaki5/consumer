package consumer

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"consumer/models"
)

func TestModule(t *testing.T) {
	t.Run("Units test of consumer", func(t *testing.T) {
		t.Run("When there is panic", func(t *testing.T) {
			t.Run("When sending Batch during panic on", func(t *testing.T) {
				c := NewConsumer()
				c.SetPanic(true)
				assert.Equal(t, errors.New(FoolServer), c.ServeBatch(make([]models.Item, 3)))
			})
			t.Run("When sending Batch of len that equal buffer after panic with no waiting set off panic", func(t *testing.T) {
				c := NewConsumer()
				c.xLocker.Lock()
				c.panic()
				c.xLocker.Unlock()
				assert.Equal(t, errors.New(FoolServer), c.ServeBatch(make([]models.Item, 3)))
			})
			t.Run("When sending Batch of len that equal buffer after panic with waiting set off panic", func(t *testing.T) {
				c := NewConsumer()
				c.xLocker.Lock()
				c.panic()
				c.xLocker.Unlock()
				time.Sleep(11 * time.Second)
				assert.Equal(t, nil, c.ServeBatch(make([]models.Item, 3)))
			})
		})
		t.Run("When there is not any panic", func(t *testing.T) {
			tests := []struct {
				name   string
				batch  models.Batch
				expErr error
			}{
				{
					name:   "When sending batch of len that less then buffer",
					batch:  make([]models.Item, 3),
					expErr: nil,
				},
				{
					name:   "When sending batch of len that more then buffer",
					batch:  make([]models.Item, 6),
					expErr: errors.New(FoolServer),
				},
				{
					name:   "When sending batch of len equal buffer",
					batch:  make([]models.Item, 5),
					expErr: nil,
				},
				{
					name:   "When sending batch of nul len",
					batch:  make([]models.Item, 0),
					expErr: nil,
				},
			}
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					c := NewConsumer()
					actualErr := c.ServeBatch(test.batch)
					assert.Equal(t, test.expErr, actualErr)
				})
			}
		})
	})
}
