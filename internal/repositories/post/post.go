package postrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	postdomain "github.com/codepnw/blog-api/internal/domains/post"
)

type postModel struct {
	ID         string    `db:"id"`
	AuthorID   string    `db:"author_id"`
	Title      string    `db:"title"`
	Content    string    `db:"content"`
	CategoryID string    `db:"category_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Repository interface {
	Insert(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error)
	FindByID(ctx context.Context, id string) (*postdomain.Post, error)
	FindByAuthorID(ctx context.Context, authorID string) ([]*postdomain.Post, error)
	List(ctx context.Context) ([]*postdomain.Post, error)
	Update(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Insert(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error) {
	m := r.inputToModel(input)
	query := `
		INSERT INTO posts (author_id, title, content, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		m.AuthorID,
		m.Title,
		m.Content,
		m.CategoryID,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*postdomain.Post, error) {
	post := new(postdomain.Post)
	query := `
		SELECT id, author_id, title, content, category_id, created_at, updated_at
		WHERE id = $1 LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.AuthorID,
		&post.Title,
		&post.Content,
		&post.CategoryID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *repository) FindByAuthorID(ctx context.Context, authorID string) ([]*postdomain.Post, error) {
	query := `
		SELECT id, author_id, title, content, category_id, created_at, updated_at
		WHERE author_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*postdomain.Post
	for rows.Next() {
		p := new(postdomain.Post)
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.CategoryID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *repository) List(ctx context.Context) ([]*postdomain.Post, error) {
	query := `
		SELECT id, author_id, title, content, category_id, created_at, updated_at
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*postdomain.Post
	for rows.Next() {
		p := new(postdomain.Post)
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.CategoryID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *repository) Update(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error) {
	var (
		sb   strings.Builder
		args []any
		idx  = 1
	)

	sb.WriteString("UPDATE posts SET ")

	if input.Title != "" {
		sb.WriteString(fmt.Sprintf("title = $%d", idx))
		args = append(args, input.Title)
		idx++
	}

	if input.Content != "" {
		sb.WriteString(fmt.Sprintf("content = $%d", idx))
		args = append(args, input.Content)
		idx++
	}

	if input.CategoryID != "" {
		sb.WriteString(fmt.Sprintf("category_id = $%d", idx))
		args = append(args, input.CategoryID)
		idx++
	}

	final := fmt.Sprintf(`
		 updated_at = NOW() 
		 WHERE id = $%d
		 RETURNING id, author_id, title, content, category_id, created_at, updated_at
	`, idx)

	sb.WriteString(final)
	args = append(args, input.ID)

	query := sb.String()
	m := new(postModel)
	log.Println(query)

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&m.ID,
		&m.AuthorID,
		&m.Title,
		&m.Content,
		&m.CategoryID,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return r.modelToDomain(m), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (r *repository) inputToModel(input *postdomain.Post) *postModel {
	return &postModel{
		ID:         input.ID,
		AuthorID:   input.AuthorID,
		Title:      input.Title,
		Content:    input.Content,
		CategoryID: input.CategoryID,
		CreatedAt:  input.CreatedAt,
		UpdatedAt:  input.UpdatedAt,
	}
}

func (r *repository) modelToDomain(input *postModel) *postdomain.Post {
	return &postdomain.Post{
		ID:         input.ID,
		AuthorID:   input.AuthorID,
		Title:      input.Title,
		Content:    input.Content,
		CategoryID: input.CategoryID,
		CreatedAt:  input.CreatedAt,
		UpdatedAt:  input.UpdatedAt,
	}
}
