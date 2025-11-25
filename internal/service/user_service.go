package service

import (
	"context"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/humooo/avito-backend-trainee-2025/internal/repo"
)

type UserService struct {
	userRepo repo.UserRepository
	prRepo   repo.PRRepository
}

func NewUserService(userRepo repo.UserRepository, prRepo repo.PRRepository) *UserService {
	return &UserService{userRepo: userRepo, prRepo: prRepo}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, active bool) error {
	return s.userRepo.SetActive(ctx, userID, active)
}

func (s *UserService) ListByTeam(ctx context.Context, teamName string, activeOnly bool) ([]*models.User, error) {
	return s.userRepo.ListByTeam(ctx, teamName, activeOnly)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetReviewPRs(ctx context.Context, reviewerID string) ([]*models.PullRequest, error) {
	return s.prRepo.ListByReviewer(ctx, reviewerID)
}
