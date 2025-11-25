package repo

import (
	"context"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	ListByTeam(ctx context.Context, teamID int64, activeOnly bool) ([]*models.User, error)
	SetActive(ctx context.Context, id int64, active bool) error
	Update(ctx context.Context, user *models.User) error
}

type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	GetByID(ctx context.Context, id int64) (*models.Team, error)
	FindByName(ctx context.Context, name string) (*models.Team, error)
}

type PRRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	GetByID(ctx context.Context, id int64) (*models.PullRequest, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	AddReviewer(ctx context.Context, prID, reviewerID int64) error
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID int64) error
	ListByReviewer(ctx context.Context, reviewerID int64) ([]*models.PullRequest, error)
}
