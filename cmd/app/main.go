package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/humooo/avito-backend-trainee-2025/internal/api"
)

func main() {
	r := chi.NewRouter()

	// Временная заглушка
	handler := &apiHandler{}

	// Используем правильную функцию из сгенерированного кода
	api.HandlerFromMux(handler, r)

	log.Println("Starting server on :8080")
	log.Println("Server is running but all endpoints return 501 Not Implemented")
	http.ListenAndServe(":8080", r)
}

// Полная заглушка всех методов
type apiHandler struct{}

func (h *apiHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *apiHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *apiHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *apiHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *apiHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *apiHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func (h *apiHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
