package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	if err = pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	sqlBytes, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		log.Printf("Warning: could not read migrations/init.sql: %v", err)
	} else {
		_, err = pool.Exec(context.Background(), string(sqlBytes))
		if err != nil {
			log.Printf("Migration warning: %v", err)
		} else {
			log.Println("Migrations applied")
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

	r.Get("/stats", handler.CustomGetStats)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
