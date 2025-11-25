package service

import (
	"context"
	"fmt"

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

// Create creates team and creates/updates members.
// members: slice of models.User where ID==0 => create, else update.
func (s *TeamService) Create(ctx context.Context, name string, members []models.User) (*models.Team, error) {
	// if team exists -> error
	if existing, _ := s.teamRepo.FindByName(ctx, name); existing != nil {
		return nil, fmt.Errorf("team exists")
	}
	team := &models.Team{Name: name}
	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}
	for _, m := range members {
		m.TeamID = team.ID
		if m.ID == 0 {
			if err := s.userRepo.Create(ctx, &m); err != nil {
				return nil, err
			}
		} else {
			// try update, if not found -> create
			if err := s.userRepo.Update(ctx, &m); err != nil {
				_ = s.userRepo.Create(ctx, &m)
			}
		}
	}
	return team, nil
}

func (s *TeamService) GetByName(ctx context.Context, name string) (*models.Team, []*models.User, error) {
	team, err := s.teamRepo.FindByName(ctx, name)
	if err != nil || team == nil {
		return nil, nil, fmt.Errorf("team not found")
	}
	users, err := s.userRepo.ListByTeam(ctx, team.ID, false)
	if err != nil {
		return nil, nil, err
	}
	return team, users, nil
}

func (s *TeamService) AddUser(ctx context.Context, teamID int64, userID int64) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || u == nil {
		return fmt.Errorf("user not found")
	}
	u.TeamID = teamID
	return s.userRepo.Update(ctx, u)
}
