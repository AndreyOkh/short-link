package main

import (
	"fmt"
	"net/http"
	"short-link/configs"
	"short-link/internal/auth"
	"short-link/internal/link"
	"short-link/pkg/db"
	"short-link/pkg/middleware"
)

func main() {

	conf := configs.LoadConfig()

	dbConn := db.NewDb(conf)
	router := http.NewServeMux()

	//	Repositories
	linkRepository := link.NewLinkRepository(dbConn)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf})
	link.NewLinkHandler(router, link.LinkHandlerDeps{LinkRepository: linkRepository})
	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Server run")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
