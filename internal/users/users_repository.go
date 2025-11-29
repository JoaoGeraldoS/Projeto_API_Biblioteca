package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *Users) error {
	sql := "INSERT INTO users (name, email, password, username, role) VALUES (?, ?, ?, ?, ?)"

	result, err := r.db.ExecContext(ctx, sql, user.Name, user.Email, user.Password, user.Username, user.Role)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id

	return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]Users, error) {
	query := `SELECT id, name, email FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var getUsers []Users

	for rows.Next() {
		var createdAtRaw, updatedAtRaw time.Time
		createdAt := createdAtRaw.Format("01/01/01 15:04:05")
		updatedAt := updatedAtRaw.Format("01/01/01 15:04:05")

		var u Users

		rows.Scan(&u.ID, &u.Name, &u.Email, &createdAt, &updatedAt)

		getUsers = append(getUsers, u)
	}

	return getUsers, nil
}

func (r *UserRepository) GetById(ctx context.Context, id int64) (*Users, error) {
	query := `SELECT id, name, email, bio FROM users WHERE id = ?`

	var user Users

	var createdAtRaw, updatedAtRaw time.Time
	createdAt := createdAtRaw.Format("01/01/01 15:04:05")
	updatedAt := updatedAtRaw.Format("01/01/01 15:04:05")

	row := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Bio, &createdAt, &updatedAt,
	)
	if row != nil {
		return nil, row
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *Users) error {
	query := "UPDATE users SET name = ? bio = ? WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Bio, user.ID)
	if err != nil {
		return err
	}

	_, errRow := result.RowsAffected()
	if errRow != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserDetails(ctx context.Context, email string) (*Users, error) {
	query := "SELECT email, username, password, role FROM users WHERE email = ?"

	var user Users
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Email, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	fmt.Println(user)

	return &user, nil
}
