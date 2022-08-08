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
	x           int
	xLocker     sync.Mutex
	panicLocker sync.Mutex
	isPanic     bool
}

func (co *Consumer) setPanic(isPanic bool) {
	co.panicLocker.Lock()
	co.isPanic = isPanic
	co.panicLocker.Unlock()
}
func (co *Consumer) ServeBatch(batch Batch) error {
	co.xLocker.Lock()
	co.panicLocker.Lock()
	if co.isPanic {
		return errors.New("server is fool")
	}
	co.panicLocker.Unlock()
	defer co.xLocker.Unlock()
	if co.x < len(batch) {
		co.Panic()
		return errors.New("server is fool")
	}
	for range batch {
		co.x--
		go func() {
			time.Sleep(2 * time.Second)
			co.xLocker.Lock()
			co.x++
			defer co.xLocker.Unlock()
		}()
	}
	return nil
}

func (co *Consumer) Panic() {
	if co.isPanic {
		return
	}
	go func() {
		co.setPanic(true)
		time.Sleep(5 * time.Minute)
		co.setPanic(false)
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
		defer c.Request().Body.Close()
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
