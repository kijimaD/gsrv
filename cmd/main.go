package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/kijimaD/gsrv"
)

func main() {
	if len(os.Args) != 2 {
		panic("Usage: [docroot]\n")
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		os.Stdout.Write([]byte("> "))
		scanner.Scan()
		in := scanner.Text()
		r := strings.NewReader(in)
		os.Stdout.Write([]byte("\n"))
		gsrv.Service(r, os.Stdout, os.Args[1])
	}
}
