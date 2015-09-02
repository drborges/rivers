package main

import (
	"flag"
	"fmt"
	"github.com/drborges/rivers"
	"github.com/drborges/rivers/stream"
	"github.com/fatih/color"
	"log"
)

var host = flag.String("host", "127.0.0.1", "The host to connect to")
var port = flag.String("port", "8181", "The port where the host is listening on")

func main() {
	flag.Parse()

	toString := func(data stream.T) stream.T { return string(data.([]byte)) }
	printMessage := func(data stream.T) { log.Println(data) }
	formatMessage := func(data stream.T) stream.T {
		red := color.New(color.FgRed)
		yellow := color.New(color.FgYellow)
		return fmt.Sprintf("[%v] %v", red.SprintFunc()("Server"), yellow.SprintFunc()(data))
	}

	rivers.DebugEnabled = true
	rivers.FromSocket("tcp", *host+":"+*port).Map(toString).Map(formatMessage).Each(printMessage).Drain()
}
