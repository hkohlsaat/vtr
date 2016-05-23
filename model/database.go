// Package model provides access to the models and separates the database management
// from the business logic of adding new users for example.
// The models usually use the CRUD methods Create, Read, Update and Delete. See their
// documentation for particular information.
package model

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "vtr.db")
	if err != nil {
		log.Fatalf("failed to connect database: %v\n", err)
	}

	if !db.HasTable(&Teacher{}) {
		db.CreateTable(&Teacher{})
		db.CreateTable(&Subject{})
		db.CreateTable(&User{})
		db.CreateTable(&dbPlan{})
		db.CreateTable(&UnknownTeacher{})
		db.CreateTable(&UnknownSubject{})
	}
}
