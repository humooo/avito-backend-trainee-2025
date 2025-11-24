package service

import "github.com/humooo/avito-backend-trainee-2025/internal/repo"

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
