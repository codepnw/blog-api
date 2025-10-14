package userrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	userdomain "github.com/codepnw/blog-api/internal/domains/user"
	"github.com/codepnw/blog-api/internal/utils/errs"
)

type UserModel struct {
	ID           string    `db:"id"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Repository interface {
	Insert(ctx context.Context, input *userdomain.User) (*userdomain.User, error)
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
	List(ctx context.Context) ([]*userdomain.User, error)
	Update(ctx context.Context, input *userdomain.User) (*userdomain.User, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Insert(ctx context.Context, input *userdomain.User) (*userdomain.User, error) {
	m := r.inputToModel(input)
	query := `
		INSERT INTO users (first_name, last_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		m.FirstName,
		m.LastName,
		m.Email,
		m.PasswordHash,
		m.Role,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*userdomain.User, error) {
	m := new(UserModel)
	query := `
		SELECT id, first_name, last_name, email, role, created_at, updated_at
		FROM users WHERE id = $1 LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.FirstName,
		&m.LastName,
		&m.Email,
		&m.Role,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *repository) List(ctx context.Context) ([]*userdomain.User, error) {
	query := `
		SELECT id, first_name, last_name, email, role, created_at, updated_at
		FROM users
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*userdomain.User
	for rows.Next() {
		u := new(userdomain.User)
		err = rows.Scan(
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Role,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *repository) Update(ctx context.Context, input *userdomain.User) (*userdomain.User, error) {
	var (
		sb   strings.Builder
		args []any
		idx  = 1
	)
	sb.WriteString("UPDATE users SET ")

	if input.FirstName != "" {
		sb.WriteString(fmt.Sprintf("first_name = $%d,", idx))
		args = append(args, input.FirstName)
		idx++
	}

	if input.LastName != "" {
		sb.WriteString(fmt.Sprintf("last_name = $%d,", idx))
		args = append(args, input.LastName)
		idx++
	}

	sb.WriteString(fmt.Sprintf(`
	 	updated_at = NOW() WHERE id = $%d
		RETURNING id, first_name, last_name, email, role, created_at, updated_at
	`, idx))
	args = append(args, input.ID)

	m := new(UserModel)
	query := sb.String()
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&m.ID,
		&m.FirstName,
		&m.LastName,
		&m.Email,
		&m.Role,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	return r.modelToDomain(m), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *repository) inputToModel(input *userdomain.User) *UserModel {
	return &UserModel{
		ID:           input.ID,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.UpdatedAt,
	}
}

func (r *repository) modelToDomain(input *UserModel) *userdomain.User {
	return &userdomain.User{
		ID:           input.ID,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		Role:         input.Role,
		CreatedAt:    input.CreatedAt,
		UpdatedAt:    input.UpdatedAt,
	}
}
