package service

import (
	"context"
	"fmt"
	"math/rand"

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

func (s *PRService) Create(ctx context.Context, title string, authorID int64) (*models.PullRequest, error) {
	author, err := s.userRepo.GetById(ctx, authorID)
	if err != nil || author == nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}
	team, err := s.teamRepo.GetById(ctx, author.TeamID)
	if err != nil || team == nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}
	users, err := s.userRepo.ListByTeam(ctx, team.ID, true)
	if err != nil {
		return nil, fmt.Errorf("users not found: %w", err)
	}
	var reviewersCandidates []*models.User
	for _, user := range users {
		if user.ID != authorID {
			reviewersCandidates = append(reviewersCandidates, user)
		}
	}
	rand.Shuffle(len(reviewersCandidates), func(i, j int) {
		reviewersCandidates[i], reviewersCandidates[j] = reviewersCandidates[j], reviewersCandidates[i]
	})
	if len(reviewersCandidates) > 2 {
		reviewersCandidates = reviewersCandidates[:2]
	}
	pr := models.PullRequest{
		Title:     title,
		AuthorID:  authorID,
		Status:    "OPEN",
		Reviewers: []int64{},
	}
	err = s.prRepo.Create(ctx, &pr)
	if err != nil {
		return nil, fmt.Errorf("не получилось создать pr: %w", err)
	}
	for _, reviewer := range reviewersCandidates {
		s.prRepo.AddReviewer(ctx, pr.ID, reviewer.ID)
		pr.Reviewers = append(pr.Reviewers, reviewer.ID)
	}
	return &pr, nil
}

func (s *PRService) Merge(ctx context.Context, prId int64) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetById(ctx, prId)
	if err != nil {
		return nil, fmt.Errorf("pr not found: %w", err)
	}
	if pr.Status == "MERGED" {
		return pr, nil
	}
	err = s.prRepo.UpdateStatus(ctx, prId, "MERGED")
	if err != nil {
		return nil, err
	}
	pr.Status = "MERGED"
	return pr, nil
}

func (s *PRService) Reassign(ctx context.Context, prID, oldReviewerID int64) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetById(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("pr not found: %w", err)
	}
	if pr.Status == "MERGED" {
		return nil, fmt.Errorf("pr is merged")
	}
	found := false
	ind := -1
	for i, id := range pr.Reviewers {
		if id == oldReviewerID {
			found = true
			ind = i
		}
	}
	if !found {
		return nil, fmt.Errorf("old reviewer not found")
	}
	oldRev, err := s.userRepo.GetById(ctx, oldReviewerID)
	if err != nil {
		return nil, fmt.Errorf("old reviewer not found: %w", err)
	}
	cands, err := s.userRepo.ListByTeam(ctx, oldRev.TeamID, true)
	if err != nil {
		return nil, fmt.Errorf("candidates not found: %w", err)
	}
	options := []int64{}
	for _, cand := range cands {
		if cand.ID != oldReviewerID {
			options = append(options, cand.ID)
		}
	}
	if len(options) == 0 {
		return nil, fmt.Errorf("candidates not found: %w", err)
	}
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	newReviewerId := options[0]
	err = s.prRepo.ReplaceReviewer(ctx, prID, oldReviewerID, newReviewerId)
	if err != nil {
		return nil, fmt.Errorf("не получилось реплейс: %w", err)
	}
	pr.Reviewers[ind] = newReviewerId
	return pr, nil
}
