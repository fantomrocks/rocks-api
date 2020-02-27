package main

import (
	"fantomrocks-api/internal/common"
	"fantomrocks-api/internal/handlers"
	"fantomrocks-api/internal/repository"
	"fantomrocks-api/internal/services"
	"net/http"
)

// Fantom Rocks API daemon serves GraphQL requests and provides details about Fantom transactions
// in the Opera/XAR block chain.
func main() {
	// load config and construct the server shared environment
	cfg := common.LoadConfig()
	log := services.NewLogger(cfg)

	// create repository
	repo, err := repository.NewRepository(cfg, log)
	if err != nil {
		log.Fatalf("Can not create application data repository. Terminating!")
	}

	// setup GraphQL API handler
	http.Handle("/api", handlers.ApiHandler(cfg, repo, log))

	// show the server opening info and start the server with DefaultServeMux
	log.Infof("Welcome to Fantom Rocks API server on [%s]", cfg.BindAddr)
	log.Fatal(http.ListenAndServe(cfg.BindAddr, nil))
}
