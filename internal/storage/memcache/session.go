package memcache

import "github.com/dmartzol/hmm/internal/hmm"

type SessionMemcache map[string]*hmm.Session

func NewSessionMemcache() *SessionMemcache {
	m := make(SessionMemcache)
	return &m
}

func (m SessionMemcache) Session(token string) (*hmm.Session, bool) {
	session, ok := m[token]
	if !ok {
		return nil, false
	}
	return session, true
}

func (m SessionMemcache) AddSession(session *hmm.Session) {
	m[session.Token] = session
}

func (m SessionMemcache) SessionFromToken(token string) (*hmm.Session, bool) {
	session, ok := m.Session(token)
	if !ok {
		return nil, false
	}
	return session, true
}
