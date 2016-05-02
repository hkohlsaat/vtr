package model

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoSuchUser          = errors.New("model: user not found")
	ErrNoMatchNamePassword = errors.New("model: wrong username or password")
)

type User struct {
	ID       uint
	Name     string `gorm:"unique_index"`
	Password string
}

func UsernameTaken(name string) bool {
	var user User
	db.Where(&User{Name: name}).First(&user)
	return user != User{}
}

func (user *User) Create(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Could't generate hashed password: %v\n", err)
	}
	user.Password = string(hash)
	db.Create(user)
}

func (user *User) Read() {
	db.Where(user).First(user)
}

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

func (user *User) Update() {
	db.Save(user)
}

func (user *User) UpdatePassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Couldn't generate hashed password: %v\n", err)
	}
	user.Password = string(hash)
	db.Save(user)
}

func (user *User) Delete() {
	db.Where(user).Delete(User{})
}
