package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	if user.ID == 0 {
		err := r.pool.QueryRow(ctx,
			"INSERT INTO users (name, is_active, team_id) VALUES ($1, $2, $3) RETURNING id",
			user.Name, user.IsActive, user.TeamID).Scan(&user.ID)

		return err

	}

	_, err := r.pool.Exec(ctx,
		"INSERT INTO users (id, name, is_active, team_id) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET name=$2, is_active=$3, team_id=$4",
		user.ID, user.Name, user.IsActive, user.TeamID)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx, "SELECT id, name, is_active, team_id FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Name, &u.IsActive, &u.TeamID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *UserRepo) ListByTeam(ctx context.Context, teamID int64, activeOnly bool) ([]*models.User, error) {
	query := "SELECT id, name, is_active, team_id FROM users WHERE team_id=$1"
	args := []any{teamID}

	if activeOnly {
		query += " AND is_active = true"
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.IsActive, &u.TeamID); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) SetActive(ctx context.Context, id int64, active bool) error {
	_, err := r.pool.Exec(ctx, "UPDATE users SET is_active=$1 WHERE id=$2", active, id)
	return err
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	tag, err := r.pool.Exec(ctx,
		"UPDATE users SET name=$1, is_active=$2, team_id=$3 WHERE id=$4",
		user.Name, user.IsActive, user.TeamID, user.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
