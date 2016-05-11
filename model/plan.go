package model

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

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

type dbPlan struct {
	gorm.Model
	FileName string
}

func (plan *Plan) Create() []byte {
	b, _ := json.Marshal(*plan)

	var filename string

	var count int
	db.Model(&dbPlan{}).Count(&count)
	if count == 0 {
		filename = plan.Created.Format("02.01.2006") + "-1.json"
	} else {
		filename = LastPlanName()
		tmp := strings.Split(filename, "-")
		oldDay := tmp[0]
		tmp = strings.Split(tmp[1], ".")
		number, _ := strconv.Atoi(tmp[0])
		day := plan.Created.Format("02.01.2006")
		if oldDay == day {
			filename = day + "-" + strconv.Itoa(number+1) + ".json"
		} else {
			filename = day + "-1.json"
		}
	}
	dbplan := &dbPlan{FileName: filename}
	db.Create(&dbplan)

	ioutil.WriteFile(filename, b, 0644)
	return b
}

func LastPlanName() string {
	var readDbPlan dbPlan
	db.Order("created_at desc").First(&readDbPlan)
	return readDbPlan.FileName
}
