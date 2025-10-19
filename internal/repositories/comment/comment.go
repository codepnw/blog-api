package commentrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	commentdomain "github.com/codepnw/blog-api/internal/domains/comment"
	"github.com/codepnw/blog-api/internal/utils/errs"
)

type CommentModel struct {
	ID        int64     `db:"id"`
	PostID    string    `db:"post_id"`
	UserID    string    `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Repository interface {
	Insert(ctx context.Context, input *commentdomain.Comment) (*commentdomain.Comment, error)
	FindByID(ctx context.Context, id int64) (*commentdomain.Comment, error)
	ListByPost(ctx context.Context, postID string) ([]*commentdomain.Comment, error)
	Update(ctx context.Context, input *commentdomain.Comment) (*commentdomain.Comment, error)
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Insert(ctx context.Context, input *commentdomain.Comment) (*commentdomain.Comment, error) {
	m := r.inputToModel(input)
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, m.PostID, m.UserID, m.Content).Scan(
		&m.ID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *repository) FindByID(ctx context.Context, id int64) (*commentdomain.Comment, error) {
	comment := new(commentdomain.Comment)
	query := `
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCommentNotFound
		}
		return nil, err
	}
	return comment, nil
}

func (r *repository) ListByPost(ctx context.Context, postID string) ([]*commentdomain.Comment, error) {
	query := `
		SELECT id, post_id, user_id, content, created_at, updated_at
		FROM comments WHERE post_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*commentdomain.Comment
	for rows.Next() {
		c := new(commentdomain.Comment)
		err = rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *repository) Update(ctx context.Context, input *commentdomain.Comment) (*commentdomain.Comment, error) {
	m := r.inputToModel(input)
	query := `
		UPDATE comments SET content = $1
		WHERE id = $2 AND post_id = $3
		RETURNING id, post_id, user_id, content, created_at, updated_at
	`
	err := r.db.QueryRowContext(ctx, query, m.Content, m.ID, m.PostID).Scan(
		&m.ID,
		&m.PostID,
		&m.UserID,
		&m.Content,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCommentNotFound
		}
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errs.ErrCommentNotFound
	}
	return nil
}

func (r *repository) inputToModel(input *commentdomain.Comment) *CommentModel {
	return &CommentModel{
		ID:        input.ID,
		PostID:    input.PostID,
		UserID:    input.UserID,
		Content:   input.Content,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}

func (r *repository) modelToDomain(input *CommentModel) *commentdomain.Comment {
	return &commentdomain.Comment{
		ID:        input.ID,
		PostID:    input.PostID,
		UserID:    input.UserID,
		Content:   input.Content,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}
}
