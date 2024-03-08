package store

import (
    "database/sql"
    "github.com/nitishsaini706/stan/assignment2/models"
    _ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = sql.ErrNoRows

type Store struct {
    db *sql.DB
}

func New(db *sql.DB) *Store {
    return &Store{db: db}
}

// CreateUser inserts a new user into the database
func (s *Store) CreateUser(user models.User) error {
    _, err := s.db.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)", user.ID, user.Name, user.Email)
    return err
}

// GetUser retrieves a user by their ID from the database
func (s *Store) GetUser(id int) (models.User, error) {
    var user models.User
    err := s.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return models.User{}, err
    }
    return user, nil
}

// UpdateUser updates an existing user's details in the database
func (s *Store) UpdateUser(id int, newUser models.User) error {
    _, err := s.db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", newUser.Name, newUser.Email, id)
    return err
}

// DeleteUser removes a user from the database
func (s *Store) DeleteUser(id int) error {
    _, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
    return err
}


// need this for database running or maigration
func Migrate(db *sql.DB) error {
    query := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL
);`

    _, err := db.Exec(query)
    return err
}