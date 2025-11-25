package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/humooo/avito-backend-trainee-2025/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PRRepo struct {
	pool *pgxpool.Pool
}

func NewPRRepo(pool *pgxpool.Pool) *PRRepo {
	return &PRRepo{pool: pool}
}

func (r *PRRepo) Create(ctx context.Context, pr *models.PullRequest) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := "INSERT INTO pull_requests (title, author_id, status) VALUES ($1, $2, $3) RETURNING id"
	err = tx.QueryRow(ctx, query, pr.Title, pr.AuthorID, pr.Status).Scan(&pr.ID)
	if err != nil {
		return err
	}

	for _, revID := range pr.Reviewers {
		query = "INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)"
		_, err := tx.Exec(ctx, query, pr.ID, revID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PRRepo) GetByID(ctx context.Context, id int64) (*models.PullRequest, error) {
	pr := &models.PullRequest{}

	query := "SELECT id, title, author_id, status FROM pull_requests WHERE id=$1"
	err := r.pool.QueryRow(ctx, query, id).Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("pr not found")
	}
	if err != nil {
		return nil, err
	}

	query = "SELECT reviewer_id FROM pr_reviewers WHERE pr_id=$1"
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rID int64
		if err := rows.Scan(&rID); err != nil {
			return nil, err
		}
		pr.Reviewers = append(pr.Reviewers, rID)
	}

	return pr, nil
}

func (r *PRRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := "UPDATE pull_requests SET status=$1 WHERE ID=$2"
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *PRRepo) AddReviewer(ctx context.Context, prID, reviewerID int64) error {
	query := "INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	_, err := r.pool.Exec(ctx, query, prID, reviewerID)
	return err
}

func (r *PRRepo) ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID int64) error {
	query := "UPDATE pr_reviewers SET reviewer_id=$1 WHERE pr_id=$2 AND reviewer_id=$3"
	tag, err := r.pool.Exec(ctx, query, newReviewerID, prID, oldReviewerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reviewer not found for replacement")
	}
	return nil
}

func (r *PRRepo) ListByReviewer(ctx context.Context, reviewerID int64) ([]*models.PullRequest, error) {
	query := `
		SELECT pr.id, pr.title, pr.author_id, pr.status
		FROM pull_requests pr
		JOIN pr_reviewers rev ON pr.id = rev.pr_id
		WHERE rev.reviewer_id = $1
	`
	rows, err := r.pool.Query(ctx, query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []*models.PullRequest
	for rows.Next() {
		pr := &models.PullRequest{}
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}
