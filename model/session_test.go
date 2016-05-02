package model

import (
	"testing"
	"time"
)

var session1, session2 Session

func TestNewSession(t *testing.T) {
	user1, user2 := User{Name: "testuser1"}, User{Name: "testuser2"}
	session1, session2 = SStore.NewSession(user1), SStore.NewSession(user2)
	if session1.Id == session2.Id || time.Now().After(session1.Expiration) || time.Now().After(session2.Expiration) {
		t.Fail()
	}
}

func TestSession(t *testing.T) {
	SStore.Lock()
	session2.Expiration = time.Now()
	SStore.sessions[session2.Id] = session2
	SStore.Unlock()

	s1, ok1 := SStore.Session(session1.Id)
	s2, ok2 := SStore.Session(session2.Id)
	s3, ok3 := SStore.Session(SStore.newsid())
	if s1 != session1 || !ok1 {
		t.Error("s1 not as expected")
	} else if s2 == session2 || !(s2 == Session{}) || ok2 {
		t.Error("s2 not as expected")
	} else if !(s3 == Session{}) || ok3 {
		t.Error("s3 not as expected")
	}

	SStore.RLock()
	_, ok2 = SStore.sessions[session2.Id]
	if ok2 {
		t.Fail()
	}
	SStore.RUnlock()
}

func TestDeleteSession(t *testing.T) {
	if _, ok := SStore.Session(session1.Id); !ok {
		t.Fail()
	}
	SStore.DeleteSession(session1.Id)
	if _, ok := SStore.Session(session1.Id); ok {
		t.Fail()
	}
}
