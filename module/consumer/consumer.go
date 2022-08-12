package consumer

import (
	"errors"
	"sync"
	"time"

	"github.com/labstack/gommon/log"

	"consumer/models"
)

type Consumer struct {
	x       int
	isPanic bool
	xLocker sync.Mutex
}

func NewConsumer() *Consumer {
	return &Consumer{
		x: ConsumerBuffer,
	}
}

func (co *Consumer) ServeBatch(batch models.Batch) error {
	co.xLocker.Lock()
	defer co.xLocker.Unlock()
	if co.isPanic {
		log.Debug(FoolServer)
		return errors.New(FoolServer)
	}
	if co.x < len(batch) {
		co.panic()
		return errors.New(FoolServer)
	}
	for range batch {
		co.x--
		go func() {
			time.Sleep(ItemServeTime)
			co.IncrX()
		}()
	}
	return nil
}

func (co *Consumer) panic() {
	co.setPanic(true)
	log.Debug(FoolServer)
	go func() {
		time.Sleep(PanicDuration)
		log.Debug("server can work")
		co.SetPanic(false)
	}()
}

func (co *Consumer) SetPanic(isPanic bool) {
	co.xLocker.Lock()
	defer co.xLocker.Unlock()

	co.setPanic(isPanic)
}

func (co *Consumer) setPanic(isPanic bool) {
	co.isPanic = isPanic
}

func (co *Consumer) IncrX() {
	co.xLocker.Lock()
	defer co.xLocker.Unlock()

	co.incrX()
}

func (co *Consumer) incrX() {
	co.x++
}
