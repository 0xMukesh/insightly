package main

import (
	"github.com/0xmukesh/insightly/cmd"
	"github.com/0xmukesh/insightly/internal/utils"
)

func main() {
	if err := cmd.Execute(); err != nil {
		utils.LogF(err.Error())
	}
}
