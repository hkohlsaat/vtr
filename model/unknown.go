package model

type UnknownTeacher struct {
	Short string `gorm:"unique_index"`
}

func (unknownTeacher *UnknownTeacher) Create() {
	if db.NewRecord(*unknownTeacher) {
		db.Create(unknownTeacher)
	}
}

func ReadAllUnknownTeachers() []UnknownTeacher {
	var unknownTeachers []UnknownTeacher
	db.Find(&unknownTeachers)
	return unknownTeachers
}

func (unknownTeacher *UnknownTeacher) Delete() {
	db.Where(unknownTeacher).Delete(UnknownTeacher{})
}

type UnknownSubject struct {
	Short string `gorm:"unique_index"`
}

func (unknownUnknownSubject *UnknownSubject) Create() {
	if db.NewRecord(*unknownUnknownSubject) {
		db.Create(unknownUnknownSubject)
	}
}

func ReadAllUnknownSubjects() []UnknownSubject {
	var unknownUnknownSubjects []UnknownSubject
	db.Find(&unknownUnknownSubjects)
	return unknownUnknownSubjects
}

func (unknownUnknownSubject *UnknownSubject) Delete() {
	db.Where(unknownUnknownSubject).Delete(UnknownSubject{})
}
