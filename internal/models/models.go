package models

import "time"

type User struct {
	ID       string
	Name     string
	IsActive bool
	TeamName string
}

type Team struct {
	Name string
}

type PullRequest struct {
	ID        string
	Title     string
	AuthorID  string
	Status    string
	Reviewers []string
	CreatedAt time.Time
	MergedAt  *time.Time
}

type UserStat struct {
	Username    string
	ReviewCount int
}
