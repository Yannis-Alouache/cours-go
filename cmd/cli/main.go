package main

import (
	"log"

	"cours-go/internal/client"
	"cours-go/internal/local"
	"cours-go/internal/tui"
)

func main() {
	cfg, err := local.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	state, err := local.LoadState()
	if err != nil {
		log.Fatal(err)
	}

	apiClient := client.New(cfg.APIURL)
	apiClient.SetToken(state.Token)

	app := tui.New(cfg, state, apiClient)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
