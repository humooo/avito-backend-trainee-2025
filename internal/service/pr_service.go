package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/humooo/avito-backend-trainee-2025/internal/repo"
)

type PRService struct {
	prRepo   repo.PRRepository
	userRepo repo.UserRepository
	teamRepo repo.TeamRepository
}

func NewPRService(prRepo repo.PRRepository, userRepo repo.UserRepository, teamRepo repo.TeamRepository) *PRService {
	return &PRService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo}
}

func (s *PRService) Create(ctx context.Context, id, title, authorID string) (*models.PullRequest, error) {
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, fmt.Errorf("author not found")
	}

	pr := &models.PullRequest{
		ID:       id,
		Title:    title,
		AuthorID: authorID,
		Status:   "OPEN",
	}

	if err := s.prRepo.CreateWithReviewers(ctx, pr); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, fmt.Errorf("pr exists")
		}
		return nil, err
	}
	return pr, nil
}

func (s *PRService) Merge(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return pr, nil
	}
	if err := s.prRepo.Merge(ctx, prID); err != nil {
		return nil, err
	}

	now := time.Now()
	pr.Status = "MERGED"
	pr.MergedAt = &now

	return pr, nil
}

func (s *PRService) Reassign(ctx context.Context, prID, oldReviewerID string) (*models.PullRequest, string, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("pr not found")
	}
	if pr.Status == "MERGED" {
		return nil, "", fmt.Errorf("pr merged")
	}

	isAssigned := false
	for _, r := range pr.Reviewers {
		if r == oldReviewerID {
			isAssigned = true
			break
		}
	}
	if !isAssigned {
		return nil, "", fmt.Errorf("reviewer not assigned")
	}

	oldUser, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil || oldUser == nil {
		return nil, "", fmt.Errorf("old reviewer not found")
	}

	newID, err := s.prRepo.FindCandidateForReassign(ctx, oldUser.TeamName, oldReviewerID, pr.AuthorID, pr.Reviewers)
	if err != nil {
		return nil, "", err
	}
	if newID == "" {
		return nil, "", fmt.Errorf("no candidates")
	}

	if err := s.prRepo.ReplaceReviewer(ctx, prID, oldReviewerID, newID); err != nil {
		return nil, "", err
	}

	for i, r := range pr.Reviewers {
		if r == oldReviewerID {
			pr.Reviewers[i] = newID
			break
		}
	}

	return pr, newID, nil
}
