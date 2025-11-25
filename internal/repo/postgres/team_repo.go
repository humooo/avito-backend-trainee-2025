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
	if team.ID == 0 {
		err := r.pool.QueryRow(ctx, "INSERT INTO teams (name) VALUES ($1) RETURNING id", team.Name).Scan(&team.ID)
		return err
	}
	_, err := r.pool.Exec(ctx, "INSERT INTO teams (id, name) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING", team.ID, team.Name)
	return err
}

func (r *TeamRepo) GetByID(ctx context.Context, id int64) (*models.Team, error) {
	t := &models.Team{}
	err := r.pool.QueryRow(ctx, "SELECT id, name FROM teams WHERE id=$1", id).Scan(&t.ID, &t.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return t, err
}

func (r *TeamRepo) FindByName(ctx context.Context, name string) (*models.Team, error) {
	t := &models.Team{}
	err := r.pool.QueryRow(ctx, "SELECT id, name FROM teams WHERE name=$1", name).Scan(&t.ID, &t.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return t, err
}
