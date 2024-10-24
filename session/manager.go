package session

import (
	"errors"
	"net/http"
)

// Replace this variable to change the key
var SessionIDKey = "session-id"

type Manager struct {
	path   string
	maxAge int
}

func NewManager(path string, maxAge int) *Manager {
	return &Manager{
		path:   path,
		maxAge: maxAge,
	}
}

func (manager Manager) Get(r *http.Request) (*Session, error) {
	session, err := r.Cookie(SessionIDKey)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return nil, err
	}

	if errors.Is(err, http.ErrNoCookie) {
		return &Session{}, nil
	}

	return &Session{id: session.Value}, nil
}

func (manager Manager) Save(w http.ResponseWriter, session *Session) {
	cookie := http.Cookie{
		Name:     SessionIDKey,
		Value:    session.id,
		MaxAge:   manager.maxAge,
		Path:     manager.path,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}
