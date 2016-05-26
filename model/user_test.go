package model

import "testing"

var defaultUser = User{Name: "testuser"}
var password = "passwd"

func TestUserCreate(t *testing.T) {
	user := defaultUser
	user.Create(password)
	if user == defaultUser {
		t.Error("user not written")
	}
	if user.Password == password {
		t.Error("raw password saved!")
	}
}

func TestUsernameTaken(t *testing.T) {
	u := User{Name: "Georg"}
	if !defaultUser.Exists() || u.Exists() {
		t.Error("taken usernames not correctly recognised")
	}
}

func TestCountUsers(t *testing.T) {
	if CountUsers() != 1 {
		t.Error("didn't count users as expected")
	}
}

func TestUserRead(t *testing.T) {
	user := defaultUser
	user.Read()
	if user == defaultUser {
		t.Error("didn't read user")
	}
}

func TestUserGetWithPassword(t *testing.T) {
	user := defaultUser
	err := user.GetWithPassword(password)
	if err != nil || user == defaultUser {
		t.Error("didn't read with password")
	}
	user = defaultUser
	err = user.GetWithPassword("123456")
	if err != ErrNoMatchNamePassword {
		t.Error("didn't recognise user password mismatch")
	}
	user = User{Name: "wrongname"}
	err = user.GetWithPassword(password)
	if err != ErrNoSuchUser {
		t.Error("didn't recognise unknown user")
	}
}

func TestUserUpdate(t *testing.T) {
	user := defaultUser
	user.Read()
	user.Name = "differentName"
	user.Update()
	user = defaultUser
	user.Read()
	if user != defaultUser {
		t.Errorf("user found by old name %+v", user)
	}
	defaultUser = User{Name: "differentName"}
	user = defaultUser
	user.Read()
	if user == defaultUser {
		t.Error("can't find user after username update")
	}
}

func TestUserUpdatePassword(t *testing.T) {
	user := defaultUser
	user.Read()
	oldUser := user
	user.UpdatePassword("asdf")
	user = defaultUser
	user.Read()
	if oldUser == user {
		t.Error("user not updated")
	}
	user = defaultUser
	err := user.GetWithPassword("asdf")
	if err != nil || user == defaultUser {
		t.Error("user not found with new password")
	}
}

func TestUserDelete(t *testing.T) {
	user := defaultUser
	user.Read()
	user.Delete()
	user = defaultUser
	user.Read()
	if user != defaultUser {
		t.Error("user not deleted")
	}
}
