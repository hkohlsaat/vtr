package model

import "testing"

// Subject dummies.
var s = []Subject{
	Subject{Short: "Md", Name: "Melden", SplitClass: false},
	Subject{Short: "Lm", Name: "Lärmen", SplitClass: true},
	Subject{Short: "Zl", Name: "Zählen", SplitClass: false}}

func TestSubjectCreate(t *testing.T) {
	// Create a subject.
	subject := s[0]
	subject.Create()

	// Check existence.
	var count int
	if db.Get(&count, "SELECT count(*) FROM subjects WHERE short = ?", s[0].Short); count == 0 {
		t.Error("Subject wasn't created.")
	}
}

func TestSubjectRead(t *testing.T) {
	// Read subject by short.
	subject := Subject{Short: s[0].Short}
	subject.Read()

	if subject.Name != s[0].Name {
		t.Errorf("Subject was not read as expected: (%+v) actual: %+v", s[0], subject)
	}

	// Read another subject by short. This one doesn't exists.
	subject = Subject{Short: s[1].Short}
	subject.Read()

	if subject.Short != s[1].Short || subject.Name != "" || subject.SplitClass != false {
		t.Errorf("Subject was unexpectedly read: %+v", subject)
	}
}

func TestSubjectExists(t *testing.T) {
	// Test that subject exists.
	subject := Subject{Short: s[0].Short}
	if !subject.Exists() {
		t.Error("Existent subject was not recognized.")
	}

	// Test that subject exists not.
	subject = Subject{Short: s[1].Short}
	if subject.Exists() {
		t.Error("Non existent subject was recognized.")
	}
}

func TestSubjectUpdate(t *testing.T) {
	// Change Name and SplitClass.
	newName := s[1].Name
	subject := Subject{Short: s[0].Short, Name: newName, SplitClass: s[0].SplitClass}
	subject.Update()

	// Test that subject.
	subject = Subject{Short: s[0].Short}
	subject.Read()
	if subject.Name != newName || subject.SplitClass != s[0].SplitClass {
		t.Error("Subject was read with old values after update.")
	}
}

func TestSubjectUpdateShort(t *testing.T) {
	// Change Short, Name and SplitClass.
	subject := s[1]
	subject.UpdateShort(s[0].Short)

	// Read old short. There shouldn't be anything to read.
	subject = Subject{Short: s[0].Short}
	if subject.Exists() {
		t.Error("Subject is still associated with old short after UpdateShort call.")
	}

	// Read updated subject.
	subject = Subject{Short: s[1].Short}
	subject.Read()
	if subject != s[1] {
		t.Error("Subject was read with old values after update.")
	}
}

func TestReadAllSubjects(t *testing.T) {
	// Count before.
	lenBefore := len(ReadAllSubjects())

	// Create new records.
	subjects := []Subject{
		Subject{Short: "t1", Name: "Test1", SplitClass: false},
		Subject{Short: "t2", Name: "Test2", SplitClass: true}}
	subjects[0].Create()
	subjects[1].Create()

	// Read all subjects and test whether the newly created subjects are returned, too.
	var hasFirst, hasSecond bool
	var readSubjects = ReadAllSubjects()
	for _, subject := range readSubjects {
		if subject == subjects[0] {
			hasFirst = true
		}
		if subject == subjects[1] {
			hasSecond = true
		}
	}
	if len(readSubjects) != lenBefore+2 {
		t.Error("The read amount of subjects differs from the expected amount.")
	}
	if !hasFirst || !hasSecond {
		t.Error("Didn't read all subjects.")
	}

	// Delete the records.
	subjects[0].Delete()
	subjects[1].Delete()
}

func TestSubjectDelete(t *testing.T) {
	// Create new record.
	subject := s[2]
	subject.Create()

	// Delete the subject.
	subject = s[1]
	subject.Delete()

	// Check if it was deleted.
	subject = Subject{Short: s[1].Short}
	if subject.Exists() {
		t.Error("Subject still exists after deletion.")
	}

	// Delete last subject.
	subject = s[2]
	subject.Delete()

	// Check if it was deleted.
	subject = Subject{Short: s[2].Short}
	if subject.Exists() {
		t.Error("Subject still exists after deletion.")
	}
}
