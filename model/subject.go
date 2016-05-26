package model

// Subject represents a subject associating abbreviations (short)
// with other subject information.
type Subject struct {
	Short      string
	Name       string
	SplitClass bool
}

const subject_schema = `CREATE TABLE subjects (short text, name text, splitclass boolean)`

// ReadAllSubjects fetches all subject records from the database and
// returns a slice with all subjects found.
func ReadAllSubjects() []Subject {
	var subjects []Subject
	db.Select(&subjects, `SELECT short, name, splitclass FROM subjects`)

	return subjects
}

// Exists tells whether there is a subject record with this subject's short.
func (s *Subject) Exists() bool {
	var count int
	db.Get(&count, "SELECT count(*) FROM subjects WHERE short = ?", s.Short)
	return count > 0
}

// Create inserts this subject into the database. This only happens
// if there isn't an entry with this subject's short already. Otherwise
// nothing happens.
func (s *Subject) Create() {
	// Verify that this subject isn't in the database already.
	if !s.Exists() {
		// This subject (subject with this short) is not in the database
		// already, so it is inserted now.
		stmt := `INSERT INTO subjects(short, name, splitclass) VALUES (?, ?, ?)`
		db.Exec(stmt, s.Short, s.Name, s.SplitClass)
	}
}

// Read completes this subject with the subject information associated
// with this subject's short.
func (s *Subject) Read() {
	db.Get(s, "SELECT short, name, splitclass FROM subjects WHERE short = ?", s.Short)
}

// Update updates the subject record with the same short as this subject's short
// with the new data. To change the short itself, use UpdateShort.
func (s *Subject) Update() {
	stmt := `UPDATE subjects SET name = ?, splitclass = ? WHERE short = ?`
	db.Exec(stmt, s.Name, s.SplitClass, s.Short)
}

// UpdateShort updates the subject identified by the given short with the
// data included in the given subject receiver.
func (s *Subject) UpdateShort(short string) {
	stmt := `UPDATE subjects SET short = ?, name = ?, splitclass = ? WHERE short = ?`
	db.Exec(stmt, s.Short, s.Name, s.SplitClass, short)
}

// Delete removes this subject from the database.
func (s *Subject) Delete() {
	stmt := `DELETE FROM subjects WHERE short = ?`
	db.Exec(stmt, s.Short)
}
