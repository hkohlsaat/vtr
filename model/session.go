package model

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"sync"
	"time"
)

var SStore *SessionStore

func init() {
	d, err := time.ParseDuration("720h")
	if err != nil {
		log.Fatal(err)
	}
	SStore = &SessionStore{CookieName: "vtr_gsp_session", MaxAge: int(d.Seconds()), sessions: make(map[string]Session), duration: d}
}

type SessionStore struct {
	// CookieName is the name of the cookie containing the session's id.
	CookieName string
	// MaxAge is the duration the session is valid in seconds to be set with the cookie.
	MaxAge int

	sync.RWMutex
	sessions map[string]Session
	duration time.Duration
}

type Session struct {
	Id         string
	Username   string
	Expiration time.Time
}

func (store *SessionStore) NewSession(user User) Session {
	store.Lock()
	session := Session{Id: store.newsid(), Username: user.Name, Expiration: time.Now().Add(store.duration)}
	store.sessions[session.Id] = session
	store.Unlock()
	return session
}

func (store *SessionStore) Session(sid string) (Session, bool) {
	store.RLock()
	session, ok := store.sessions[sid]
	store.RUnlock()

	// Check whether the store had got the session.
	if ok {
		// Check whether the session is exipired or valid.
		if session.Expiration.After(time.Now()) {
			return session, true
		} else {
			store.DeleteSession(session.Id)
		}
	}
	// No session found.
	return Session{}, false
}

func (store *SessionStore) DeleteSession(sid string) {
	store.Lock()
	_, ok := store.sessions[sid]
	if ok {
		delete(store.sessions, sid)
	}
	store.Unlock()
}

func (store *SessionStore) newsid() string {
	var sid string
	for {
		b := make([]byte, 32)
		_, err := io.ReadFull(rand.Reader, b)
		if err != nil {
			log.Fatal(err)
		}
		sid = base64.URLEncoding.EncodeToString(b)
		if _, ok := store.sessions[sid]; !ok {
			break
		}
	}
	return sid
}
