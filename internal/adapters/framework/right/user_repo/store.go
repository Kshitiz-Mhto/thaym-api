package user_repo

import (
	"database/sql"
	"fmt"

	"ecom-api/internal/application/core/types/entity"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user entity.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password, isVerified) VALUES (?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password, user.IsVerified)
	if err != nil {
		return fmt.Errorf("failed to retreive users: %w", err)
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*entity.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by email: %w", err)
	}

	defer rows.Close()

	user := new(entity.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// if u.ID == 0 {
	// 	log.Printf("noooooooooooooooo9")

	// 	return nil, fmt.Errorf("user not found")
	// }

	return user, nil
}

func (s *Store) GetUserByID(id string) (*entity.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by id: %w", err)
	}

	defer rows.Close()

	user := new(entity.User)
	for rows.Next() {
		user, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUsersByRole(role string) ([]*entity.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE role = ?", role)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by role: %w", err)
	}

	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		user, err := scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Store) SetUserLocking(userEmail string, isLocked bool) error {
	_, err := s.db.Exec("UPDATE users SET isLocked = ? WHERE email = ?", isLocked, userEmail)

	if err != nil {
		return fmt.Errorf("failed to lock user: %w", err)
	}
	return nil
}

func scanRowsIntoUser(rows *sql.Rows) (*entity.User, error) {
	user := new(entity.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsVerified,
		&user.Role,
		&user.IsLocked,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
