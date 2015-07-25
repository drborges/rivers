package main

import (
	"github.com/drborges/riversv2"
	"net"
	"log"
)

func main() {
	stream := rivers.New().FromSocket("tcp", ":8081").Sink()

	ln, _ := net.Listen("tcp", ":8081")
	go func() {
		conn, _ := ln.Accept()
		defer conn.Close()
		conn.Write([]byte("Hello\n"))
		conn.Write([]byte("there\n"))
		conn.Write([]byte("rivers!\n"))
	}()

	log.Println(stream.Read())
}