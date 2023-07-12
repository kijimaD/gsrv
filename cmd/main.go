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
		scanner.Scan()
		in := scanner.Text()
		r := strings.NewReader(in)
		gsrv.Service(r, os.Stdout, os.Args[1])
	}
}
