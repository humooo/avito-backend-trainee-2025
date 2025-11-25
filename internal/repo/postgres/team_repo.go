package postgres

import (
	"context"
	"errors"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

func (r *TeamRepo) Create(ctx context.Context, team *models.Team) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO teams (name) VALUES ($1)", team.Name)
	return err
}

func (r *TeamRepo) FindByName(ctx context.Context, name string) (*models.Team, error) {
	t := &models.Team{}
	err := r.pool.QueryRow(ctx, "SELECT name FROM teams WHERE name=$1", name).Scan(&t.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return t, err
}
