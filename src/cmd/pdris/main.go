package main

import (
	"log"

	"gihub.com/bongerka/sberPDRIS/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	a.Run()
}
