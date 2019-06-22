package main

import (
	"log"

	"github.com/copybird/copybird/common"
)

func main() {
	app := common.NewApp()
	if err := app.Run(); err != nil {
		log.Printf("run err: %s", err)
	}
}
