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
	x        int
	mxLocker sync.Mutex
	isPanic  bool
}

func (co *Consumer) ServeBatch(batch Batch) error {
	co.mxLocker.Lock()
	defer co.mxLocker.Unlock()
	if co.x < len(batch) {
		co.Panic()
		return errors.New("server is fool")
	}
	var wg sync.WaitGroup
	for range batch {
		co.x--
		wg.Add(1)
		go func() {
			time.Sleep(2 * time.Second)
			co.x++
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func (co *Consumer) Panic() {
	if co.isPanic {
		return
	}
	go func(sm *Consumer) {
		sm.isPanic = true
		time.Sleep(5 * time.Minute)
		sm.isPanic = false
	}(co)

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

		return c.String(http.StatusOK, "ok")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
