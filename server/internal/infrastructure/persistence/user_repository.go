package persistence

import (
	"database/sql"
	"fmt"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/database"
	"webscraper-v2/pkg/datetime"
)

type userRepository struct {
	db *database.SQLiteDB
}

func NewUserRepository(db *database.SQLiteDB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	query := `INSERT INTO users (username, email, password, role, active, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	res, err := r.db.Exec(query,
		user.Username, user.Email, user.Password, user.Role,
		user.Active, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	id, err := res.LastInsertId()

	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	user.ID = id
	return nil
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	query := `SELECT id, username, email, password, role, active, created_at, updated_at 
			  FROM users WHERE username = ? AND active = true`

	user := &entity.User{}
	var createdAt, updatedAt string
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.Active, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by username: %w", err)
	}
	user.CreatedAt, err = datetime.Parse(createdAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}
	user.UpdatedAt, err = datetime.Parse(updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %w", err)
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	query := `SELECT id, username, email, password, role, active, created_at, updated_at 
			  FROM users WHERE email = ? AND active = true`

	user := &entity.User{}
	var createdAt, updatedAt string
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.Active, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}
	user.CreatedAt, err = datetime.Parse(createdAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}
	user.UpdatedAt, err = datetime.Parse(updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %w", err)
	}
	return user, nil
}

func (r *userRepository) FindByID(id int64) (*entity.User, error) {
	query := `SELECT id, username, email, password, role, active, created_at, updated_at 
			  FROM users WHERE id = ? AND active = true`

	user := &entity.User{}
	var createdAt, updatedAt string
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.Active, &createdAt, &updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by id: %w", err)
	}
	user.CreatedAt, err = datetime.Parse(createdAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %w", err)
	}
	user.UpdatedAt, err = datetime.Parse(updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error parsing updated_at: %w", err)
	}
	return user, nil
}

func (r *userRepository) Update(user *entity.User) error {
	query := `UPDATE users SET username = ?, email = ?, password = ?, role = ?, active = ?, updated_at = ? 
			  WHERE id = ?`

	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		user.Username, user.Email, user.Password, user.Role,
		user.Active, user.UpdatedAt, user.ID)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(id int64) error {
	query := `UPDATE users SET active = false, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, time.Now(), id)

	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

func (r *userRepository) ExistsUsername(username string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE username = ? AND active = true`
	var count int
	err := r.db.QueryRow(query, username).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking username existence: %w", err)
	}
	return count > 0, nil
}

func (r *userRepository) ExistsEmail(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ? AND active = true`
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking email existence: %w", err)
	}
	return count > 0, nil
}
