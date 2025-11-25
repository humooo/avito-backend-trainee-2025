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

func (r *PRRepo) CreateWithReviewers(ctx context.Context, pr *models.PullRequest) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO pull_requests (id, title, author_id, status) VALUES ($1, $2, $3, $4)",
		pr.ID, pr.Title, pr.AuthorID, pr.Status)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO pr_reviewers (pr_id, reviewer_id)
		SELECT $1, u.id
		FROM users u
		JOIN users author ON author.id = $2
		WHERE u.team_name = author.team_name
		  AND u.id != $2
		  AND u.is_active = TRUE
		ORDER BY RANDOM()
		LIMIT 2
		RETURNING reviewer_id
	`
	rows, err := tx.Query(ctx, query, pr.ID, pr.AuthorID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var rID string
		if err := rows.Scan(&rID); err != nil {
			return err
		}
		pr.Reviewers = append(pr.Reviewers, rID)
	}

	return tx.Commit(ctx)
}

func (r *PRRepo) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	pr := &models.PullRequest{}
	err := r.pool.QueryRow(ctx, "SELECT id, title, author_id, status, created_at, merged_at FROM pull_requests WHERE id=$1", id).
		Scan(&pr.ID, &pr.Title, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("pr not found")
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, "SELECT reviewer_id FROM pr_reviewers WHERE pr_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rID string
		rows.Scan(&rID)
		pr.Reviewers = append(pr.Reviewers, rID)
	}
	return pr, nil
}

func (r *PRRepo) Merge(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, "UPDATE pull_requests SET status='MERGED', merged_at=NOW() WHERE id=$1", id)
	return err
}

func (r *PRRepo) ReplaceReviewer(ctx context.Context, prID, oldID, newID string) error {
	tag, err := r.pool.Exec(ctx, "UPDATE pr_reviewers SET reviewer_id=$1 WHERE pr_id=$2 AND reviewer_id=$3", newID, prID, oldID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reviewer not found")
	}
	return nil
}

func (r *PRRepo) ListByReviewer(ctx context.Context, reviewerID string) ([]*models.PullRequest, error) {
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

func (r *PRRepo) FindCandidateForReassign(ctx context.Context, teamName, oldReviewerID, authorID string, currentReviewers []string) (string, error) {
	query := `
		SELECT id FROM users
		WHERE team_name=$1 AND is_active=TRUE AND id!=$2 AND id!=$3 AND id != ALL($4)
		ORDER BY RANDOM() LIMIT 1
	`
	var newID string
	err := r.pool.QueryRow(ctx, query, teamName, oldReviewerID, authorID, currentReviewers).Scan(&newID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	return newID, err
}
