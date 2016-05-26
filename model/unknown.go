package model

type UnknownTeacher struct {
	Short string
}

const unknown_schema = "CREATE TABLE unknown_teachers (short TEXT UNIQUE);CREATE TABLE unknown_subjects (short TEXT UNIQUE)"

func (ut *UnknownTeacher) Create() {
	stmt := `INSERT INTO unknown_teachers (short) VALUES (?)`
	db.Exec(stmt, ut.Short)
}

func ReadAllUnknownTeachers() []UnknownTeacher {
	var unknownTeachers []UnknownTeacher
	db.Select(&unknownTeachers, "SELECT short FROM unknown_teachers")
	return unknownTeachers
}

func (ut *UnknownTeacher) Delete() {
	stmt := `DELETE FROM unknown_teachers WHERE short = ?`
	db.Exec(stmt, ut.Short)
}

type UnknownSubject struct {
	Short string
}

func (us *UnknownSubject) Create() {
	stmt := `INSERT INTO unknown_subjects (short) VALUES (?)`
	db.Exec(stmt, us.Short)
}

func ReadAllUnknownSubjects() []UnknownSubject {
	var unknownSubjects []UnknownSubject
	db.Select(&unknownSubjects, "SELECT short FROM unknown_subjects")
	return unknownSubjects
}

func (us *UnknownSubject) Delete() {
	stmt := `DELETE FROM unknown_subjects WHERE short = ?`
	db.Exec(stmt, us.Short)
}
