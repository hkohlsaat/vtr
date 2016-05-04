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

// User struct conveys user data, primarily the user's name.
type User struct {
	ID       uint
	Name     string `gorm:"unique_index"`
	Password string
}

// UsernameTaken returns true if the username in name is already in use and therefore
// can't be used by a second user.
func UsernameTaken(name string) bool {
	var user User
	db.Where(&User{Name: name}).First(&user)
	return user != User{}
}

// CountUsers returns the total of registered users in the database.
func CountUsers() int {
	var count int
	db.Model(&User{}).Count(&count)
	return count
}

// Create inserts this new user into the database.
// *user should only convey the name of this new user.
func (user *User) Create(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Could't generate hashed password: %v\n", err)
	}
	user.Password = string(hash)
	db.Create(user)
}

// Read makes user to point to the user with the name user.Name.
func (user *User) Read() {
	db.Where(user).First(user)
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
func (user *User) GetWithPassword(password string) error {
	var readUser User
	db.Where(user).First(&readUser)
	if readUser.Name == user.Name {
		err := bcrypt.CompareHashAndPassword([]byte(readUser.Password), []byte(password))
		switch {
		case err == nil:
			*user = readUser
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
func (user *User) Update() {
	db.Save(user)
}

// UpdatePassword saves the new password (hash) into the existing user's entry.
// This method hashes the password properly.
func (user *User) UpdatePassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Couldn't generate hashed password: %v\n", err)
	}
	user.Password = string(hash)
	db.Save(user)
}

// Delete removes the user from the database.
func (user *User) Delete() {
	db.Where(user).Delete(User{})
}
