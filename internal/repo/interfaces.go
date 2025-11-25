package repo

import (
	"context"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
)

type UserRepository interface {
	Upsert(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	ListByTeam(ctx context.Context, teamName string, activeOnly bool) ([]*models.User, error)
	SetActive(ctx context.Context, id string, active bool) error
}

type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	FindByName(ctx context.Context, name string) (*models.Team, error)
}

type PRRepository interface {
	CreateWithReviewers(ctx context.Context, pr *models.PullRequest) error
	GetByID(ctx context.Context, id string) (*models.PullRequest, error)
	Merge(ctx context.Context, id string) error
	ReplaceReviewer(ctx context.Context, prID, oldID, newID string) error
	ListByReviewer(ctx context.Context, reviewerID string) ([]*models.PullRequest, error)
	FindCandidateForReassign(ctx context.Context, teamName, oldReviewerID, authorID string, currentReviewers []string) (string, error)
}
