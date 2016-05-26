package model

import (
	"encoding/json"
	"time"
)

// Plan represents the substitution's plan file with both parts.
type Plan struct {
	Created time.Time
	Parts   []Part
}

// Part represents the list of substitutions for one day.
type Part struct {
	Day           time.Time
	Substitutions []Substitution
}

// Substitutions represents the substitution's information.
type Substitution struct {
	Period       string
	Class        string
	SubstTeacher Teacher
	InstdTeacher Teacher
	InstdSubject Subject
	Kind         string
	Text         string
	TaskProvider Teacher
}

const plan_schema = `CREATE TABLE plans (upload DATETIME UNIQUE, json TEXT, file BLOB)`

// Create saves this plan to the database as the newest plan.
func (plan *Plan) Create(file []byte) {
	json, _ := json.Marshal(*plan)
	upload := time.Now()

	stmt := `INSERT INTO plans (upload, json, file) VALUES (?, ?, ?)`
	db.Exec(stmt, upload, string(json), file)
}

// LastPlanJSON returns the last plan in JSON format.
func LastPlanJSON() string {
	var json string
	db.Get(&json, "SELECT json FROM plans ORDER BY upload DESC LIMIT 1")
	return json
}
