package main

import (
	"github.com/drborges/riversv2"
	"github.com/drborges/riversv2/rx"
	"log"
	"fmt"
	"github.com/fatih/color"
	"flag"
)

var host = flag.String("host", "127.0.0.1", "The host to connect to")
var port = flag.String("port", "8181", "The port where the host is listening on")

func main() {
	flag.Parse()

	toString := func(data rx.T) rx.T { return string(data.([]byte)) }
	printMessage := func(data rx.T) { log.Println(data) }
	formatMessage := func(data rx.T) rx.T {
		red := color.New(color.FgRed)
		yellow := color.New(color.FgYellow)
		return fmt.Sprintf("[%v] %v", red.SprintFunc()("Server"), yellow.SprintFunc()(data))
	}

	rivers.New().FromSocket("tcp", *host + ":" + *port).Map(toString).Map(formatMessage).Each(printMessage).Drain()
}
