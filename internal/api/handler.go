package api

import (
	"encoding/json"
	"net/http"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/humooo/avito-backend-trainee-2025/internal/service"
)

type ApiHandler struct {
	PRService   *service.PRService
	UserService *service.UserService
	TeamService *service.TeamService
}

// Хелпер для отправки ошибок в формате generated ErrorResponse
func (h *ApiHandler) writeError(w http.ResponseWriter, code ErrorResponseErrorCode, message string, status int) {
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func (h *ApiHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var body PostTeamAddJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, NOTFOUND, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	var members []models.User
	for _, m := range body.Members {
		members = append(members, models.User{
			ID:       m.UserId,
			Name:     m.Username,
			IsActive: m.IsActive,
		})
	}

	team, err := h.TeamService.Create(r.Context(), body.TeamName, members)
	if err != nil {
		if err.Error() == "team exists" {
			h.writeError(w, TEAMEXISTS, "team already exists", http.StatusBadRequest)
			return
		}
		h.writeError(w, NOTFOUND, err.Error(), http.StatusInternalServerError)
		return
	}

	users, _ := h.UserService.ListByTeam(r.Context(), team.Name, false)

	resp := mapTeamToResponse(team, users)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ApiHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params GetTeamGetParams) {
	team, users, err := h.TeamService.GetByName(r.Context(), params.TeamName)
	if err != nil {
		h.writeError(w, NOTFOUND, "team not found", http.StatusNotFound)
		return
	}

	resp := mapTeamToResponse(team, users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ApiHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var body PostUsersSetIsActiveJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, NOTFOUND, "invalid json", http.StatusBadRequest)
		return
	}

	err := h.UserService.SetIsActive(r.Context(), body.UserId, body.IsActive)
	if err != nil {
		h.writeError(w, NOTFOUND, "user not found", http.StatusNotFound)
		return
	}

	u, _ := h.UserService.GetByID(r.Context(), body.UserId)

	resp := User{
		UserId:   u.ID,
		Username: u.Name,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]User{"user": resp})
}

func (h *ApiHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var body PostPullRequestCreateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, NOTFOUND, "invalid json", http.StatusBadRequest)
		return
	}

	pr, err := h.PRService.Create(r.Context(), body.PullRequestId, body.PullRequestName, body.AuthorId)
	if err != nil {
		if err.Error() == "pr exists" {
			h.writeError(w, PREXISTS, "pr id already exists", http.StatusConflict)
			return
		}
		if err.Error() == "author not found" || err.Error() == "team not found" {
			h.writeError(w, NOTFOUND, err.Error(), http.StatusNotFound)
			return
		}
		h.writeError(w, NOTFOUND, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]PullRequest{"pr": mapPRToResponse(pr)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *ApiHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var body PostPullRequestMergeJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, NOTFOUND, "invalid json", http.StatusBadRequest)
		return
	}

	pr, err := h.PRService.Merge(r.Context(), body.PullRequestId)
	if err != nil {
		h.writeError(w, NOTFOUND, "pr not found", http.StatusNotFound)
		return
	}

	resp := map[string]PullRequest{"pr": mapPRToResponse(pr)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ApiHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var body PostPullRequestReassignJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.writeError(w, NOTFOUND, "invalid json", http.StatusBadRequest)
		return
	}

	pr, newID, err := h.PRService.Reassign(r.Context(), body.PullRequestId, body.OldUserId)
	if err != nil {
		switch err.Error() {
		case "pr merged":
			h.writeError(w, PRMERGED, "cannot reassign on merged PR", http.StatusConflict)
		case "reviewer not assigned":
			h.writeError(w, NOTASSIGNED, "reviewer is not assigned to this PR", http.StatusConflict)
		case "no candidates":
			h.writeError(w, NOCANDIDATE, "no active replacement candidate in team", http.StatusConflict)
		case "pr not found", "old reviewer not found":
			h.writeError(w, NOTFOUND, err.Error(), http.StatusNotFound)
		default:
			h.writeError(w, NOTFOUND, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := struct {
		PR         PullRequest `json:"pr"`
		ReplacedBy string      `json:"replaced_by"`
	}{
		PR:         mapPRToResponse(pr),
		ReplacedBy: newID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ApiHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params GetUsersGetReviewParams) {
	prs, err := h.UserService.GetReviewPRs(r.Context(), params.UserId)
	if err != nil {
		h.writeError(w, NOTFOUND, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]PullRequestShort, 0)
	for _, p := range prs {
		out = append(out, mapPullRequestShort(p))
	}

	response := struct {
		UserId       string             `json:"user_id"`
		PullRequests []PullRequestShort `json:"pull_requests"`
	}{
		UserId:       params.UserId,
		PullRequests: out,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
