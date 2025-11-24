package models

type User struct {
	ID       int64
	Name     string
	IsActive bool
	TeamID   int64
}

type Team struct {
	ID   int64
	Name string
}

type PullRequest struct {
	ID        int64
	Title     string
	AuthorID  int64
	Status    string
	Reviewers []int64
}
