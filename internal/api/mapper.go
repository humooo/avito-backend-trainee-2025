package api

import (
	"fmt"
	"strconv"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
)

// PRToResponse converts internal PullRequest -> api.PullRequest (types.gen.go)
func PRToResponse(pr *models.PullRequest) PullRequest {
	resp := PullRequest{
		PullRequestId:     fmt.Sprintf("%d", pr.ID),
		PullRequestName:   pr.Title,
		AuthorId:          fmt.Sprintf("%d", pr.AuthorID),
		Status:            PullRequestStatus(pr.Status),
		AssignedReviewers: []string{},
	}
	for _, r := range pr.Reviewers {
		resp.AssignedReviewers = append(resp.AssignedReviewers, fmt.Sprintf("%d", r))
	}
	return resp
}

// PRToShort converts internal PullRequest -> api.PullRequestShort
func PRToShort(pr *models.PullRequest) PullRequestShort {
	return PullRequestShort{
		AuthorId:        fmt.Sprintf("%d", pr.AuthorID),
		PullRequestId:   fmt.Sprintf("%d", pr.ID),
		PullRequestName: pr.Title,
		Status:          PullRequestShortStatus(pr.Status),
	}
}

// TeamMembersFromModels converts []*models.User -> []TeamMember
func TeamMembersFromModels(users []*models.User) []TeamMember {
	var out []TeamMember
	for _, u := range users {
		out = append(out, TeamMember{
			UserId:   strconv.FormatInt(u.ID, 10),
			Username: u.Name,
			IsActive: u.IsActive,
		})
	}
	return out
}

// TeamToResponse converts internal team + users -> api.Team
func TeamToResponse(team *models.Team, users []*models.User) Team {
	return Team{
		TeamName: team.Name,
		Members:  TeamMembersFromModels(users),
	}
}
