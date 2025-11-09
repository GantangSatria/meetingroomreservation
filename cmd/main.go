package main

import (
	"log"
	"meetingroomreservation/config"
	"meetingroomreservation/internal/bootstrap"
)

func main() {
	cfg := config.LoadConfig()
	app := bootstrap.NewApp(cfg)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
