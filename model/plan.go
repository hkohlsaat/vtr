package model

import "time"

type Plan struct {
	Created time.Time
	Parts   []Part
}

type Part struct {
	Day           time.Time
	Substitutions []Substitution
}

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
