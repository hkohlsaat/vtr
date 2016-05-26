package model

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Errors returned by GetWithPassword.
// They should be used to tell what didn't work in the authentification.
var (
	ErrNoSuchUser          = errors.New("model: user not found")
	ErrNoMatchNamePassword = errors.New("model: wrong username or password")
)

// User represents a user associating his identification and authentification information.
type User struct {
	ID       uint
	Name     string
	Password string
}

const user_schema = `CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT UNIQUE, password TEXT)`

// CountUsers returns the total of recorded users in the database.
func CountUsers() int {
	var count int
	db.Get(&count, "SELECT count(*) FROM users")
	return count
}

// UsernameTaken returns true if the username is already in use and therefore
// can't be taken by a second user.
func (u *User) Exists() bool {
	var count int
	db.Get(&count, "SELECT count(*) FROM users WHERE name = ?", u.Name)
	return count > 0
}

// Create inserts this new user into the database.
// The user receiver should only convey the name of this new user.
func (u *User) Create(password string) {
	if !u.Exists() {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Couldn't generate hashed password: %v", err)
		}
		u.Password = string(hash)

		stmt := `INSERT INTO users(name, password) VALUES (?, ?)`
		db.Exec(stmt, u.Name, u.Password)
	}
}

// Read completes the user with the information associated with this user's name.
func (u *User) Read() {
	db.Get(u, "SELECT id, password FROM users WHERE name = ?", u.Name)
}

// GetWithPassword finds out whether the given password matches the password of user.
// If these don't match an error is returned. This can be ErrNoMatchNamePassword, ErrNoSuchUser
// or an error indicating that something in the password hashing went wrong.
//	err := user.GetWithPassword(password)
//	switch {
//	case err == nil:
//		// The user can get logged in or somthing.
//	case err == model.ErrNoMatchNamePassword:
//		// Maybe return a message to the user.
//	case err == model.ErrNoSuchUser:
//		// Maybe tell the user the username was wrong.
//	default:
//		// Some kind of internal error.
//	}
func (u *User) GetWithPassword(password string) error {
	var user User
	db.Get(&user, "SELECT id, name, password FROM users WHERE name = ?", u.Name)

	if user.Name == u.Name {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		switch {
		case err == nil:
			*u = user
			return nil
		case err == bcrypt.ErrMismatchedHashAndPassword:
			return ErrNoMatchNamePassword
		default:
			return err
		}
	} else {
		return ErrNoSuchUser
	}
}

// Update saves this user (changes the information in the existing entry).
// The user should have been read before (with the Read method).
// For updating the password use UpdatePassword.
func (u *User) Update() {
	stmt := `UPDATE users SET name = ? WHERE id = ?`
	db.Exec(stmt, u.Name, u.ID)
}

// UpdatePassword saves the new password (hash) into the existing user's entry.
// This method hashes the password properly.
func (u *User) UpdatePassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Couldn't generate hashed password: %v\n", err)
	}

	u.Password = string(hash)

	stmt := `UPDATE users SET name = ?, password = ? WHERE id = ?`
	db.Exec(stmt, u.Name, u.Password, u.ID)
}

// Delete removes this user from the database.
func (u *User) Delete() {
	stmt := `DELETE FROM users WHERE id = ?`
	db.Exec(stmt, u.ID)
}
