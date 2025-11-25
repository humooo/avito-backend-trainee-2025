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

func (r *UserRepo) Upsert(ctx context.Context, user *models.User) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO users (id, username, is_active, team_name) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET username=$2, is_active=$3, team_name=$4",
		user.ID, user.Name, user.IsActive, user.TeamName)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*models.User, error) {
	u := &models.User{}
	err := r.pool.QueryRow(ctx, "SELECT id, username, is_active, team_name FROM users WHERE id=$1", id).
		Scan(&u.ID, &u.Name, &u.IsActive, &u.TeamName)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

func (r *UserRepo) ListByTeam(ctx context.Context, teamName string, activeOnly bool) ([]*models.User, error) {
	query := "SELECT id, username, is_active, team_name FROM users WHERE team_name=$1"
	args := []any{teamName}

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
		if err := rows.Scan(&u.ID, &u.Name, &u.IsActive, &u.TeamName); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) SetActive(ctx context.Context, id string, active bool) error {
	cmd, err := r.pool.Exec(ctx, "UPDATE users SET is_active=$1 WHERE id=$2", active, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

type UserStat struct {
	Username    string `json:"username"`
	ReviewCount int    `json:"review_count"`
}

func (r *UserRepo) GetStats(ctx context.Context) ([]models.UserStat, error) {
	query := `
		SELECT u.username, COUNT(r.pr_id)
		FROM users u
		LEFT JOIN pr_reviewers r ON u.id = r.reviewer_id
		GROUP BY u.id, u.username
		ORDER BY COUNT(r.pr_id) DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.UserStat
	for rows.Next() {
		var s models.UserStat
		if err := rows.Scan(&s.Username, &s.ReviewCount); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
