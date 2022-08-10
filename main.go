package main

import (
	"consumer/module/consumer"
	"consumer/routing"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	co := consumer.NewConsumer()

	e.POST("/", routing.NewHandler(co).PostBatch)

	e.Logger.Fatal(e.Start(":1323"))
}
