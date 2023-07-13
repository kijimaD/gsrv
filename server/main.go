package main

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/kijimaD/gsrv"
)

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
		go handler(conn)
	}
}

func handler(conn *net.TCPConn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	gsrv.Service(r, conn, ".")
	time.Sleep(time.Second)
}
