package main

import (
	"log"

	"github.com/say8hi/addrforge/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
