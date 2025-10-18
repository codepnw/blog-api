package categoryrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	categorydomain "github.com/codepnw/blog-api/internal/domains/category"
)

type Repository interface {
	Insert(ctx context.Context, input *categorydomain.Category) error
	FindByID(ctx context.Context, id string) (*categorydomain.Category, error)
	List(ctx context.Context) ([]*categorydomain.Category, error)
	Update(ctx context.Context, input *categorydomain.Category) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Insert(ctx context.Context, input *categorydomain.Category) error {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, input.Name, input.Description)
	return err
}

func (r *repository) FindByID(ctx context.Context, id string) (*categorydomain.Category, error) {
	c := new(categorydomain.Category)

	query := `
		SELECT id, name, description
		FROM categories WHERE id = $1 LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Description,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *repository) List(ctx context.Context) ([]*categorydomain.Category, error) {
	query := `SELECT id, name, description FROM categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*categorydomain.Category
	for rows.Next() {
		c := new(categorydomain.Category)
		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func (r *repository) Update(ctx context.Context, input *categorydomain.Category) error {
	var (
		sb   strings.Builder
		args []any
		idx  = 1
	)
	sb.WriteString("UPDATE categories SET ")

	if input.Name != "" {
		sb.WriteString(fmt.Sprintf("name = $%d,", idx))
		args = append(args, input.Name)
		idx++
	}

	if input.Description != "" {
		sb.WriteString(fmt.Sprintf(" description = $%d", idx))
		args = append(args, input.Description)
		idx++
	}

	sb.WriteString(fmt.Sprintf(" WHERE id = $%d", idx))
	args = append(args, input.ID)

	query := sb.String()

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category not found")
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("category not found")
	}
	return nil
}
