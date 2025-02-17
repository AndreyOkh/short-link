package main

import (
	"fmt"
	"net/http"
	"short-link/configs"
	"short-link/internal/auth"
	"short-link/internal/link"
	"short-link/internal/stat"
	"short-link/internal/user"
	"short-link/pkg/db"
	"short-link/pkg/event"
	"short-link/pkg/middleware"
)

func main() {

	conf := configs.LoadConfig()

	dbConn := db.NewDb(conf)
	router := http.NewServeMux()

	eventBus := event.NewEventBus()

	//	Repositories
	linkRepository := link.NewLinkRepository(dbConn)
	userRepository := user.NewUserRepository(dbConn)
	statRepository := stat.NewStatRepository(dbConn)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.ServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config:         conf,
		EventBus:       eventBus,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	go statService.AddClick()
	fmt.Println("Server run")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
