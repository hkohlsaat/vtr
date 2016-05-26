// Package model provides access to the models and separates the database management
// from the business logic of adding new users for example.
// The models usually use the CRUD methods Create, Read, Update and Delete. See their
// documentation for particular information.
package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func init() {
	db = sqlx.MustOpen("sqlite3", "file:vtr.db?cache=shared&mode=rwc")

	// Read all table names.
	tables := tables()

	// Create tables which are not yet created.
	if !tables["teachers"] {
		db.MustExec(teacher_schema)
		db.MustExec(subject_schema)
		db.MustExec(user_schema)
		db.MustExec(plan_schema)
	}
}

func tables() map[string]bool {
	// Query all table names.
	var names []string
	db.Select(&names, `SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)

	// Fill all table names into a map
	used := make(map[string]bool)
	for _, name := range names {
		used[name] = true
	}

	return used
}
