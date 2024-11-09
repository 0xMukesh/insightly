package main

import (
	"log"

	"github.com/0xmukesh/ratemywebsite/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
