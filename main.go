package main

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Item struct{}

type Batch []Item

type Consumer struct {
	x       int
	isPanic bool
	xLocker sync.Mutex
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

func (co *Consumer) ServeBatch(batch Batch) error {
	co.xLocker.Lock()
	defer co.xLocker.Unlock()
	if co.isPanic {
		return errors.New("server is fool")
	}
	if co.x < len(batch) {
		co.panic()
		return errors.New("server is fool")
	}
	for range batch {
		co.x--
		go func() {
			time.Sleep(2 * time.Second)
			co.IncrX()
		}()
	}
	return nil
}

func (co *Consumer) panic() {
	co.setPanic(true)
	go func() {
		time.Sleep(10 * time.Second)
		log.Println("Server can work")

		co.SetPanic(false)
	}()
}

func main() {
	e := echo.New()
	var batch Batch
	var co = Consumer{
		x: 5,
	}
	e.POST("/", func(c echo.Context) error {
		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Println("Cannot read body")
			return c.String(http.StatusInternalServerError, "Cannot read body")
		}
		defer func() {
			if err = c.Request().Body.Close(); err != nil {
				log.Printf(err.Error())
			}
		}()
		err = json.Unmarshal(b, &batch)
		if err != nil {
			log.Println("Cannot unmarshal json")
			return c.String(http.StatusInternalServerError, "Cannot unmarshal json")
		}
		err = co.ServeBatch(batch)

		if err != nil {
			log.Println("Server is full")
			return c.String(http.StatusInternalServerError, "Server is full")
		}

		return c.String(http.StatusOK, "Ok")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
