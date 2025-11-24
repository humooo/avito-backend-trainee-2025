package memory

import (
	"context"
	"sync"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
)

type MemoryTeamRepo struct {
	mu     sync.Mutex
	teams  map[int64]*models.Team
	nextID int64
}

func NewMemoryTeamRepo() *MemoryTeamRepo {
	return &MemoryTeamRepo{
		teams:  make(map[int64]*models.Team),
		nextID: 1,
	}
}

func (r *MemoryTeamRepo) Create(ctx context.Context, team *models.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	team.ID = r.nextID
	r.nextID++
	r.teams[team.ID] = team
	return nil
}

func (r *MemoryTeamRepo) GetById(ctx context.Context, id int64) (*models.Team, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.teams[id], nil
}
