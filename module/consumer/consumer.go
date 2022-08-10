package consumer

import (
	"consumer/models"
	"errors"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	x       int
	isPanic bool
	xLocker sync.Mutex
}

func NewConsumer() *Consumer {
	return &Consumer{
		x:       ConsumerBuffer,
		isPanic: false,
	}
}

func (co *Consumer) ServeBatch(batch models.Batch) error {
	co.xLocker.Lock()
	defer co.xLocker.Unlock()
	if co.isPanic {
		log.Println(FoolServer)
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
	go func() {
		time.Sleep(PanicDuration)
		log.Println("Server can work")
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
