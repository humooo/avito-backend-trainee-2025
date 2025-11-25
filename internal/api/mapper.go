package api

import (
	"github.com/humooo/avito-backend-trainee-2025/internal/models"
)

func mapPRToResponse(pr *models.PullRequest) PullRequest {
	status := PullRequestStatusOPEN
	if pr.Status == "MERGED" {
		status = PullRequestStatusMERGED
	}

	return PullRequest{
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Title,
		AuthorId:          pr.AuthorID,
		Status:            status,
		AssignedReviewers: pr.Reviewers,
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func mapPullRequestShort(pr *models.PullRequest) PullRequestShort {
	status := PullRequestShortStatusOPEN
	if pr.Status == "MERGED" {
		status = PullRequestShortStatusMERGED
	}
	return PullRequestShort{
		PullRequestId:   pr.ID,
		PullRequestName: pr.Title,
		AuthorId:        pr.AuthorID,
		Status:          status,
	}
}

func mapTeamToResponse(team *models.Team, users []*models.User) Team {
	members := make([]TeamMember, len(users))
	for i, u := range users {
		members[i] = TeamMember{
			UserId:   u.ID,
			Username: u.Name,
			IsActive: u.IsActive,
		}
	}
	return Team{
		TeamName: team.Name,
		Members:  members,
	}
}
