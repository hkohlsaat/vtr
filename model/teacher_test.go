package model

import "testing"

func TestTeacherCreate(t *testing.T) {
	teacher := &Teacher{Short: "Wd", Name: "Wittdorf", Sex: 'w'}
	teacher.Create()
	teacher = &Teacher{Short: "Wd"}
	teacher.Read()
	if teacher.Name != "Wittdorf" || teacher.Sex != 'w' {
		t.Fail()
	}
}

func TestTeacherRead(t *testing.T) {
	teacher := &Teacher{Short: "Wd"}
	teacher.Read()
	if teacher.Name != "Wittdorf" {
		t.Fail()
	}
	teacher = &Teacher{Short: "Ar"}
	teacher.Read()
	if teacher.Short != "Ar" || teacher.Name != "" || teacher.Sex != 0 {
		t.Fail()
	}
}

func TestTeacherUpdate(t *testing.T) {
	teacher := &Teacher{Short: "Wd", Name: "Wedenbruck", Sex: 'm'}
	teacher.Update()
	teacher = &Teacher{Short: "Wd"}
	teacher.Read()
	if teacher.Name != "Wedenbruck" || teacher.Sex != 'm' {
		t.Fail()
	}
}

func TestReadAllTeachers(t *testing.T) {
	teachers := []*Teacher{&Teacher{Short: "t1", Name: "Test1", Sex: 'm'}, &Teacher{Short: "t2", Name: "Test2", Sex: 'm'}}
	teachers[0].Create()
	teachers[1].Create()

	allTeachers := ReadAllTeachers()
	if len(allTeachers) != 3 {
		t.Error("didn't read three teachers")
	}

	teachers[0].Delete()
	teachers[1].Delete()
}

func TestTeacherDelete(t *testing.T) {
	teacher := &Teacher{Short: "Ar", Name: "Armann", Sex: 'm'}
	teacher.Create()
	teacher = &Teacher{Short: "Wd"}
	teacher.Delete()
	teacher = &Teacher{Short: "Wd"}
	teacher.Read()
	if teacher.Short != "Wd" || teacher.Name != "" || teacher.Sex != 0 {
		t.Fail()
	}
	teacher = &Teacher{Short: "Ar"}
	teacher.Read()
	if teacher.Short != "Ar" || teacher.Name != "Armann" || teacher.Sex != 'm' {
		t.Fail()
	}
	teacher.Delete()
}
