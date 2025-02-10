package main

import (
	"fmt"
	"net/http"
	"short-link/internal/auth"
	"short-link/internal/hello"
)

func main() {

	// conf := configs.LoadConfig()

	router := http.NewServeMux()
	hello.NewHelloHandler(router)
	auth.NewAuthHandler(router)

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server run")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
