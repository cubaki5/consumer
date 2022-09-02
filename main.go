package main

import (
	"consumer/module/consumer"
	"consumer/routing"
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	debug := flag.Bool("debug", true, "sets log level to debug")

	flag.Parse()
	log.SetLevel(log.ERROR)
	if *debug {
		log.SetLevel(log.DEBUG)
	}

	e := echo.New()
	co := consumer.NewConsumer()
	e.POST("/", routing.NewHandler(co).PostBatch)

	e.Logger.Fatal(e.Start(":1323"))
}
