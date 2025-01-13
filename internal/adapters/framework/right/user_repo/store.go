package user_repo

import (
	"database/sql"

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
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*entity.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
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

func (s *Store) GetUserByID(id int) (*entity.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
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
