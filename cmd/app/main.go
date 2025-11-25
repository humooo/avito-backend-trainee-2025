package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/humooo/avito-backend-trainee-2025/internal/api"
	"github.com/humooo/avito-backend-trainee-2025/internal/repo/postgres"
	"github.com/humooo/avito-backend-trainee-2025/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/avito_db?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	sqlBytes, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		log.Printf("Warning: could not read migrations/init.sql: %v", err)
	} else {
		_, err = pool.Exec(context.Background(), string(sqlBytes))
		if err != nil {
			log.Printf("Warning: migration failed (maybe already exists): %v", err)
		} else {
			log.Println("Migrations applied successfully")
		}
	}

	r := chi.NewRouter()

	userRepo := postgres.NewUserRepo(pool)
	teamRepo := postgres.NewTeamRepo(pool)
	prRepo := postgres.NewPRRepo(pool)

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
	http.ListenAndServe(":8080", r)
}
