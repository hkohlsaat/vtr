package model

import "testing"

func TestSubjectCreate(t *testing.T) {
	subject := &Subject{Short: "Zk", Name: "Zauberkunst", SplitClass: false}
	subject.Create()
	subject = &Subject{Short: "Zk"}
	subject.Read()
	if subject.Name != "Zauberkunst" || subject.SplitClass != false {
		t.Fail()
	}
}

func TestSubjectRead(t *testing.T) {
	subject := &Subject{Short: "Zk"}
	subject.Read()
	if subject.Name != "Zauberkunst" {
		t.Fail()
	}
	subject = &Subject{Short: "Yt"}
	subject.Read()
	if subject.Short != "Yt" || subject.Name != "" || subject.SplitClass != false {
		t.Fail()
	}
}

func TestSubjectExists(t *testing.T) {
	subject := &Subject{Short: "Zk"}
	if !subject.Exists() {
		t.Error("didn't recognise existing subject")
	}
	subject = &Subject{Short: "No"}
	if subject.Exists() {
		t.Error("did recognise non existing subject")
	}
}

func TestSubjectUpdate(t *testing.T) {
	subject := &Subject{Short: "Zk", Name: "Zehnkampf", SplitClass: true}
	subject.Update()
	subject = &Subject{Short: "Zk"}
	subject.Read()
	if subject.Name != "Zehnkampf" || subject.SplitClass != true {
		t.Fail()
	}
}

func TestSubjectUpdateShort(t *testing.T) {
	subject := &Subject{Short: "Hh", Name: "Hundehutte", SplitClass: false}
	subject.UpdateShort("Zk")
	subject = &Subject{Short: "Zk"}
	subject.Read()
	if subject.Name != "" {
		t.Error("didn't update short")
	}
	subject = &Subject{Short: "Hh"}
	subject.Read()
	if subject.Name != "Hundehutte" || (subject.SplitClass != false) {
		t.Errorf("didn't update subject correctly: %+v", *subject)
	}
}

func TestReadAllSubjects(t *testing.T) {
	subjects := []*Subject{&Subject{Short: "t1", Name: "Test1", SplitClass: false}, &Subject{Short: "t2", Name: "Test2", SplitClass: false}}
	subjects[0].Create()
	subjects[1].Create()

	allSubjects := ReadAllSubjects()
	if len(allSubjects) != 3 {
		t.Error("didn't read three subjects")
	}

	subjects[0].Delete()
	subjects[1].Delete()
}

func TestSubjectDelete(t *testing.T) {
	subject := &Subject{Short: "Yt", Name: "Yetikunde", SplitClass: true}
	subject.Create()
	subject = &Subject{Short: "Zk"}
	subject.Delete()
	subject = &Subject{Short: "Zk"}
	subject.Read()
	if subject.Short != "Zk" || subject.Name != "" || subject.SplitClass != false {
		t.Fail()
	}
	subject = &Subject{Short: "Yt"}
	subject.Read()
	if subject.Short != "Yt" || subject.Name != "Yetikunde" || subject.SplitClass != true {
		t.Fail()
	}
	subject.Delete()
}
