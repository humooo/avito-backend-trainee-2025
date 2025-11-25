package service

import (
	"context"
	"fmt"

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

func (s *UserService) SetIsActive(ctx context.Context, userID int64, active bool) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || u == nil {
		return fmt.Errorf("user not found")
	}
	return s.userRepo.SetActive(ctx, userID, active)
}

func (s *UserService) ListByTeam(ctx context.Context, teamID int64, activeOnly bool) ([]*models.User, error) {
	return s.userRepo.ListByTeam(ctx, teamID, activeOnly)
}

func (s *UserService) GetReviewPRs(ctx context.Context, reviewerID int64) ([]*models.PullRequest, error) {
	return s.prRepo.ListByReviewer(ctx, reviewerID)
}
