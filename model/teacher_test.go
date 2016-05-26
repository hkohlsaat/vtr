package model

import "testing"

// Teacher dummies.
var teacherDummies = []Teacher{
	Teacher{Short: "Md", Name: "Meden", Sex: "m"},
	Teacher{Short: "Lm", Name: "Lehmann", Sex: "w"},
	Teacher{Short: "Zl", Name: "Zommerland", Sex: "m"}}

func TestTeacherCreate(t *testing.T) {
	// Create a teacher.
	teacher := teacherDummies[0]
	teacher.Create()

	// Check existence.
	var count int
	if db.Get(&count, "SELECT count(*) FROM teachers WHERE short = ?", teacherDummies[0].Short); count == 0 {
		t.Error("Teacher wasn't created.")
	}
}

func TestTeacherRead(t *testing.T) {
	// Read teacher by short.
	teacher := Teacher{Short: teacherDummies[0].Short}
	teacher.Read()

	if teacher.Name != teacherDummies[0].Name {
		t.Errorf("Teacher was not read as expected: (%+v) actual: %+v", teacherDummies[0], teacher)
	}

	// Read another teacher by short. This one doesn't exists.
	teacher = Teacher{Short: teacherDummies[1].Short}
	teacher.Read()

	if teacher.Short != teacherDummies[1].Short || teacher.Name != "" || teacher.Sex != "" {
		t.Errorf("Teacher was unexpectedly read: %+v", teacher)
	}
}

func TestTeacherExists(t *testing.T) {
	// Test that teacher exists.
	teacher := Teacher{Short: teacherDummies[0].Short}
	if !teacher.Exists() {
		t.Error("Existent teacher was not recognized.")
	}

	// Test that teacher exists not.
	teacher = Teacher{Short: teacherDummies[1].Short}
	if teacher.Exists() {
		t.Error("Non existent teacher was recognized.")
	}
}

func TestTeacherUpdate(t *testing.T) {
	// Change Name and Sex.
	newName := teacherDummies[1].Name
	teacher := Teacher{Short: teacherDummies[0].Short, Name: newName, Sex: teacherDummies[0].Sex}
	teacher.Update()

	// Test that teacher.
	teacher = Teacher{Short: teacherDummies[0].Short}
	teacher.Read()
	if teacher.Name != newName || teacher.Sex != teacherDummies[0].Sex {
		t.Error("Teacher was read with old values after update.")
	}
}

func TestTeacherUpdateShort(t *testing.T) {
	// Change Short, Name and Sex.
	teacher := teacherDummies[1]
	teacher.UpdateShort(teacherDummies[0].Short)

	// Read old short. There shouldn't be anything to read.
	teacher = Teacher{Short: teacherDummies[0].Short}
	if teacher.Exists() {
		t.Error("Teacher is still associated with old short after UpdateShort call.")
	}

	// Read updated teacher.
	teacher = Teacher{Short: teacherDummies[1].Short}
	teacher.Read()
	if teacher != teacherDummies[1] {
		t.Error("Teacher was read with old values after update.")
	}
}

func TestReadAllTeachers(t *testing.T) {
	// Count before.
	lenBefore := len(ReadAllTeachers())

	// Create new records.
	teachers := []Teacher{
		Teacher{Short: "t1", Name: "Test1", Sex: "m"},
		Teacher{Short: "t2", Name: "Test2", Sex: "w"}}
	teachers[0].Create()
	teachers[1].Create()

	// Read all teachers and test whether the newly created teachers are returned, too.
	var hasFirst, hasSecond bool
	var readTeachers = ReadAllTeachers()
	for _, teacher := range readTeachers {
		if teacher == teachers[0] {
			hasFirst = true
		}
		if teacher == teachers[1] {
			hasSecond = true
		}
	}
	if len(readTeachers) != lenBefore+2 {
		t.Error("The read amount of teachers differs from the expected amount.")
	}
	if !hasFirst || !hasSecond {
		t.Error("Didn't read all teachers.")
	}

	// Delete the records.
	teachers[0].Delete()
	teachers[1].Delete()
}

func TestTeacherDelete(t *testing.T) {
	// Create new record.
	teacher := teacherDummies[2]
	teacher.Create()

	// Delete the teacher.
	teacher = teacherDummies[1]
	teacher.Delete()

	// Check if it was deleted.
	teacher = Teacher{Short: teacherDummies[1].Short}
	if teacher.Exists() {
		t.Error("Teacher still exists after deletion.")
	}

	// Delete last teacher.
	teacher = teacherDummies[2]
	teacher.Delete()

	// Check if it was deleted.
	teacher = Teacher{Short: teacherDummies[2].Short}
	if teacher.Exists() {
		t.Error("Teacher still exists after deletion.")
	}
}
