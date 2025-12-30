package main

import (
	"droplet/internal/command"
	"log"
	"os"
)

func main() {
	app := command.NewApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
