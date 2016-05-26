package model

// Teacher represents a teacher associating abbreviations (short)
// with name and compellation information.
type Teacher struct {
	Short string
	Name  string
	Sex   string
}

const teacher_schema = `CREATE TABLE teachers (short TEXT UNIQUE, name TEXT, sex TEXT)`

// ReadAllTeachers fetches all teacher records from the database and
// returns a slice with all teachers found.
func ReadAllTeachers() []Teacher {
	var teachers []Teacher
	db.Select(&teachers, `SELECT short, name, sex FROM teachers ORDER BY name asc`)

	return teachers
}

// Exists tells whether there is a teacher record with this teacher's short.
func (t *Teacher) Exists() bool {
	var count int
	db.Get(&count, "SELECT count(*) FROM teachers WHERE short = ?", t.Short)
	return count > 0
}

// Create inserts this teacher into the database. This only happens
// if there isn't an entry with this teacher's short already. Otherwise
// nothing happens.
func (t *Teacher) Create() {
	// Verify that this teacher isn't in the database already.
	if !t.Exists() {
		// This teacher (teacher with this short) is not in the database
		// already, so it is inserted now.
		stmt := `INSERT INTO teachers(short, name, sex) VALUES (?, ?, ?)`
		db.Exec(stmt, t.Short, t.Name, t.Sex)
	}
}

// Read completes this teacher with the teacher information associated
// with this teacher's short.
func (t *Teacher) Read() {
	db.Get(t, "SELECT short, name, sex FROM teachers WHERE short = ?", t.Short)
}

// Update updates the teacher record with the same short as this teacher's short
// with the new data. To change the short itself, use UpdateShort.
func (t *Teacher) Update() {
	stmt := `UPDATE teachers SET name = ?, sex = ? WHERE short = ?`
	db.Exec(stmt, t.Name, t.Sex, t.Short)
}

// UpdateShort updates the teacher identified by the given short with the
// data included in the given teacher receiver.
func (t *Teacher) UpdateShort(short string) {
	stmt := `UPDATE teachers SET short = ?, name = ?, sex = ? WHERE short = ?`
	db.Exec(stmt, t.Short, t.Name, t.Sex, short)
}

// Delete removes this teacher from the database.
func (t *Teacher) Delete() {
	stmt := `DELETE FROM teachers WHERE short = ?`
	db.Exec(stmt, t.Short)
}
