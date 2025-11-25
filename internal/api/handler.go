package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/humooo/avito-backend-trainee-2025/internal/service"
)

// ApiHandler implements ServerInterface (from server.gen.go)
type ApiHandler struct {
	PRService   *service.PRService
	UserService *service.UserService
	TeamService *service.TeamService
}

func (h *ApiHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body PostPullRequestCreateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	authorID, err := strconv.ParseInt(body.AuthorId, 10, 64)
	if err != nil {
		http.Error(w, "invalid author_id", http.StatusBadRequest)
		return
	}
	pr, err := h.PRService.Create(ctx, body.PullRequestName, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// convert to API model
	writeJSONWithStatus(w, PRToResponse(pr), http.StatusCreated)
}

func (h *ApiHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body PostPullRequestMergeJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prID, err := strconv.ParseInt(body.PullRequestId, 10, 64)
	if err != nil {
		http.Error(w, "invalid pull_request_id", http.StatusBadRequest)
		return
	}
	pr, err := h.PRService.Merge(ctx, prID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, PRToResponse(pr))
}

func (h *ApiHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body PostPullRequestReassignJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prID, err := strconv.ParseInt(body.PullRequestId, 10, 64)
	if err != nil {
		http.Error(w, "invalid pull_request_id", http.StatusBadRequest)
		return
	}
	oldID, err := strconv.ParseInt(body.OldUserId, 10, 64)
	if err != nil {
		http.Error(w, "invalid old_user_id", http.StatusBadRequest)
		return
	}
	pr, err := h.PRService.Reassign(ctx, prID, oldID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, PRToResponse(pr))
}

func (h *ApiHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body PostTeamAddJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert api.Team.Members -> []models.User
	var members []models.User
	for _, m := range body.Members {
		var uid int64
		if m.UserId != "" {
			if v, err := strconv.ParseInt(m.UserId, 10, 64); err == nil {
				uid = v
			}
		}
		members = append(members, models.User{
			ID:       uid,
			Name:     m.Username,
			IsActive: m.IsActive,
		})
	}

	team, err := h.TeamService.Create(ctx, body.TeamName, members)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, _ := h.UserService.ListByTeam(ctx, team.ID, false)
	resp := TeamToResponse(team, users)
	writeJSONWithStatus(w, resp, http.StatusCreated)
}

func (h *ApiHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params GetTeamGetParams) {
	ctx := r.Context()
	// team_name is string per OpenAPI
	teamName := params.TeamName
	team, users, err := h.TeamService.GetByName(ctx, teamName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resp := TeamToResponse(team, users)
	writeJSON(w, resp)
}

func (h *ApiHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params GetUsersGetReviewParams) {
	ctx := r.Context()
	uid, err := strconv.ParseInt(params.UserId, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	prs, err := h.UserService.GetReviewPRs(ctx, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var out []PullRequestShort
	for _, p := range prs {
		out = append(out, PRToShort(p))
	}
	writeJSON(w, struct {
		UserId       string             `json:"user_id"`
		PullRequests []PullRequestShort `json:"pull_requests"`
	}{
		UserId:       params.UserId,
		PullRequests: out,
	})
}

func (h *ApiHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body PostUsersSetIsActiveJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(body.UserId, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	if err := h.UserService.SetIsActive(ctx, id, body.IsActive); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// simple OK response as OpenAPI does not define a full user response here
	writeJSONWithStatus(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any) {
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONWithStatus(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
