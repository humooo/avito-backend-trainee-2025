package service

import "github.com/humooo/avito-backend-trainee-2025/internal/repo"

type PRService struct {
	prRepo   repo.PRRepository
	userRepo repo.UserRepository
	teamRepo repo.TeamRepository
}

func NewPRService(prRepo repo.PRRepository, userRepo repo.UserRepository, teamRepo repo.TeamRepository) *PRService {
	return &PRService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo}
}

type UserService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

type TeamService struct {
	teamRepo repo.TeamRepository
	userRepo repo.UserRepository
}

func NewTeamService(teamRepo repo.TeamRepository, userRepo repo.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo}
}
