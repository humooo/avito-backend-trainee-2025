package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/humooo/avito-backend-trainee-2025/internal/repo"
)

type TeamService struct {
	teamRepo repo.TeamRepository
	userRepo repo.UserRepository
}

func NewTeamService(teamRepo repo.TeamRepository, userRepo repo.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo}
}

func (s *TeamService) Create(ctx context.Context, name string, members []models.User) (*models.Team, error) {
	existing, _ := s.teamRepo.FindByName(ctx, name)
	if existing != nil {
		return nil, fmt.Errorf("team exists")
	}

	team := &models.Team{Name: name}
	if err := s.teamRepo.Create(ctx, team); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, fmt.Errorf("team exists")
		}
		return nil, err
	}

	for _, m := range members {
		m.TeamName = name
		if err := s.userRepo.Upsert(ctx, &m); err != nil {
			return nil, fmt.Errorf("failed to add user %s: %w", m.ID, err)
		}
	}

	return team, nil
}

func (s *TeamService) GetByName(ctx context.Context, name string) (*models.Team, []*models.User, error) {
	team, err := s.teamRepo.FindByName(ctx, name)
	if err != nil || team == nil {
		return nil, nil, fmt.Errorf("team not found")
	}
	users, err := s.userRepo.ListByTeam(ctx, name, false)
	return team, users, err
}
