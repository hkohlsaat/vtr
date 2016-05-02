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

func TestSubjectUpdate(t *testing.T) {
	subject := &Subject{Short: "Zk", Name: "Zehnkampf", SplitClass: true}
	subject.Update()
	subject = &Subject{Short: "Zk"}
	subject.Read()
	if subject.Name != "Zehnkampf" || subject.SplitClass != true {
		t.Fail()
	}
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
