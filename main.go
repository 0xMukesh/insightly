package main

import (
	"github.com/0xmukesh/ratemywebsite/cmd"
	"github.com/0xmukesh/ratemywebsite/internal/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		utils.LogF(err.Error())
	}
}
