package main

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/kijimaD/gsrv"
)

func echoHandler(conn *net.TCPConn) {
	defer conn.Close()
	r := strings.NewReader("GET /dummy.txt HTTP/1.0")
	gsrv.Service(r, conn, ".")
	time.Sleep(time.Second)
}

var TimeoutSec = time.Second * 10

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":7777")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		err := ln.SetDeadline(time.Now().Add(TimeoutSec))
		if err != nil {
			log.Fatal(err)
		}
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go echoHandler(conn)
	}
}
