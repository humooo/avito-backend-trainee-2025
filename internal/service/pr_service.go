package service

import (
	"context"
	"fmt"
	"math/rand"
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
	rand.Seed(time.Now().UnixNano())
	return &PRService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo}
}

func (s *PRService) Create(ctx context.Context, title string, authorID int64) (*models.PullRequest, error) {
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, fmt.Errorf("author not found: %w", err)
	}
	team, err := s.teamRepo.GetByID(ctx, author.TeamID)
	if err != nil || team == nil {
		return nil, fmt.Errorf("team not found: %w", err)
	}

	users, err := s.userRepo.ListByTeam(ctx, team.ID, true)
	if err != nil {
		return nil, fmt.Errorf("users not found: %w", err)
	}

	// 1. Фильтрация (исключаем автора)
	candidates := make([]int64, 0)
	for _, u := range users {
		if u.ID != authorID {
			candidates = append(candidates, u.ID)
		}
	}

	// 2. Выбор случайных (до 2-х)
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	limit := 2
	if len(candidates) < 2 {
		limit = len(candidates)
	}
	selected := candidates[:limit]

	// 3. Создание объекта
	pr := &models.PullRequest{
		Title:    title,
		AuthorID: authorID,
		Status:   "OPEN",
		// Сразу присваиваем выбранных ревьюверов, чтобы избежать дублирования при .Create() + .AddReviewer()
		Reviewers: selected,
	}

	// 4. Сохранение
	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, fmt.Errorf("failed to create pr: %w", err)
	}

	// Цикл AddReviewer убран, так как Reviewers уже внутри pr
	return pr, nil
}

func (s *PRService) Merge(ctx context.Context, prId int64) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prId)
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
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("pr not found: %w", err)
	}
	if pr.Status == "MERGED" {
		return nil, fmt.Errorf("pr is merged")
	}

	foundIndex := -1
	currentReviewersSet := make(map[int64]bool)
	for i, id := range pr.Reviewers {
		currentReviewersSet[id] = true
		if id == oldReviewerID {
			foundIndex = i
		}
	}

	if foundIndex == -1 {
		return nil, fmt.Errorf("old reviewer not found in this PR")
	}

	oldRev, err := s.userRepo.GetByID(ctx, oldReviewerID)
	if err != nil {
		return nil, fmt.Errorf("old reviewer user not found: %w", err)
	}

	cands, err := s.userRepo.ListByTeam(ctx, oldRev.TeamID, true)
	if err != nil {
		return nil, fmt.Errorf("candidates list error: %w", err)
	}

	options := []int64{}
	for _, cand := range cands {
		// Фильтр: не старый ревьювер, не автор, и не второй текущий ревьювер
		if cand.ID != oldReviewerID && cand.ID != pr.AuthorID && !currentReviewersSet[cand.ID] {
			options = append(options, cand.ID)
		}
	}

	if len(options) == 0 {
		return nil, fmt.Errorf("no available candidates for reassign")
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	newReviewerId := options[0]

	err = s.prRepo.ReplaceReviewer(ctx, prID, oldReviewerID, newReviewerId)
	if err != nil {
		return nil, fmt.Errorf("failed to replace: %w", err)
	}

	pr.Reviewers[foundIndex] = newReviewerId
	return pr, nil
}
