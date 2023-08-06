package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/kijimaD/gsrv"
)

var TimeoutSec = time.Second * 10

func main() {
	if len(os.Args) != 2 {
		panic("Usage: [docroot]\n")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":7777")
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		go handler(conn, os.Args[1])
	}
}

func handler(conn *net.TCPConn, docroot string) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	gsrv.Service(r, conn, docroot)
	time.Sleep(100 * time.Millisecond)
}
