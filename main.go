package main

import (
	"flag"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"consumer/module/consumer"
	"consumer/routing"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()
	log.SetLevel(log.ERROR)
	if *debug {
		log.SetLevel(log.DEBUG)
	}

	e := echo.New()
	co := consumer.NewConsumer()
	strings.Split("nya", ",")
	e.POST("/", routing.NewHandler(co).PostBatch)
	e.GET("/buffer", routing.NewHandler(co).GetBuffer)
	e.Logger.Fatal(e.Start(":1323"))
}
