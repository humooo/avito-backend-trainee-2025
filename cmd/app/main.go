package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/humooo/avito-backend-trainee-2025/internal/api"
	"github.com/humooo/avito-backend-trainee-2025/internal/repo/memory"
	"github.com/humooo/avito-backend-trainee-2025/internal/service"
)

func main() {
	r := chi.NewRouter()

	userRepo := memory.NewMemoryUserRepo()
	teamRepo := memory.NewMemoryTeamRepo()
	prRepo := memory.NewMemoryPRRepo()

	userService := service.NewUserService(userRepo, prRepo)
	teamService := service.NewTeamService(teamRepo, userRepo)
	prService := service.NewPRService(prRepo, userRepo, teamRepo)

	handler := &api.ApiHandler{
		PRService:   prService,
		UserService: userService,
		TeamService: teamService,
	}
	api.HandlerFromMux(handler, r)
	log.Println("Starting server on :8080")
	log.Println("Server is running but all endpoints return 501 Not Implemented")
	http.ListenAndServe(":8080", r)
}
