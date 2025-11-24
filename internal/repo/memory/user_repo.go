package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
)

type MemoryUserRepo struct {
	mu     sync.Mutex
	users  map[int64]*models.User
	nextID int64
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{
		users:  make(map[int64]*models.User),
		nextID: 1,
	}
}

func (r *MemoryUserRepo) Create(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepo) GetById(ctx context.Context, id int64) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.users[id], nil
}

func (r *MemoryUserRepo) ListByTeam(ctx context.Context, teamID int64, activeOnly bool) ([]*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var res []*models.User
	for _, user := range r.users {
		if user.TeamID == teamID {
			if activeOnly && !user.IsActive {
				continue
			}
			res = append(res, user)
		}
	}
	return res, nil
}

func (r *MemoryUserRepo) SetActive(ctx context.Context, id int64, active bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}
	user.IsActive = active
	return nil
}
