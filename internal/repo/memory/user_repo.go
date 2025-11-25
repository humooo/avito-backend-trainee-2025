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

	// ИСПРАВЛЕНИЕ: Если ID задан извне (например, при создании команды с конкретными ID),
	// мы должны использовать его, а не перезаписывать.
	if user.ID != 0 {
		// Если вставляем ID больше текущего счетчика, подтягиваем счетчик
		if user.ID >= r.nextID {
			r.nextID = user.ID + 1
		}
	} else {
		// Иначе генерируем новый
		user.ID = r.nextID
		r.nextID++
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Возвращаем nil, если ключа нет (безопасно для логики сервиса)
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
			// Важно возвращать копию или тот же указатель (в in-memory указатель ок)
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

func (r *MemoryUserRepo) Update(ctx context.Context, user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return fmt.Errorf("not found")
	}
	r.users[user.ID] = user
	return nil
}
