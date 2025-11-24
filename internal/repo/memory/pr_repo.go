package memory

import (
	"context"
	"sync"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
)

type MemoryPRRepo struct {
	mu     sync.Mutex
	prs    map[int64]*models.PullRequest
	nextID int64
}

func NewMemoryPRRepo() *MemoryPRRepo {
	return &MemoryPRRepo{
		prs:    make(map[int64]*models.PullRequest),
		nextID: 1,
	}
}

func (r *MemoryPRRepo) Create(ctx context.Context, pr *models.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	pr.ID = r.nextID
	r.nextID++
	r.prs[pr.ID] = pr
	return nil
}

func (r *MemoryPRRepo) GetById(ctx context.Context, id int64) (*models.PullRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.prs[id], nil
}

func (r *MemoryPRRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prs[id].Status = status
	return nil
}

func (r *MemoryPRRepo) AddReviewer(ctx context.Context, prID, reviewerID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prs[prID].Reviewers = append(r.prs[prID].Reviewers, reviewerID)
	return nil
}

func (r *MemoryPRRepo) ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, id := range r.prs[prID].Reviewers {
		if id == oldReviewerID {
			r.prs[prID].Reviewers[i] = newReviewerID
		}
	}
	return nil
}

func (r *MemoryPRRepo) ListByReviewer(ctx context.Context, reviewerID int64) ([]*models.PullRequest, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var res []*models.PullRequest
	for _, pr := range r.prs {
		for _, id := range pr.Reviewers {
			if id == reviewerID {
				res = append(res, pr)
				break
			}
		}
	}
	return res, nil
}
