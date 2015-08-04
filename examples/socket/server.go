package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"github.com/fatih/color"
	"flag"
)

var host = flag.String("host", "0.0.0.0", "The host to connect to")
var port = flag.String("port", "8181", "The port where the host is listening on")

func main() {
	flag.Parse()

	yellow := color.New(color.FgYellow)
	ln, _ := net.Listen("tcp", *host + ":" + *port)
	log.Println("Listening at:", yellow.SprintFunc()(ln.Addr()))

	reader := bufio.NewReader(os.Stdin)
	conn, _ := ln.Accept()

	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		conn.Write(bytes)
	}
}
