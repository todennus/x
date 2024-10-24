package session

import (
	"context"

	"github.com/gorilla/securecookie"
	"github.com/todennus/x/xreflect"
)

var _ Store[int] = (*CookieStore[int])(nil)

type CookieStore[T any] struct {
	codecs []securecookie.Codec
}

func NewCookieStore[T any](authenticationKey, encryptionKey []byte) *CookieStore[T] {
	return &CookieStore[T]{
		codecs: securecookie.CodecsFromPairs(authenticationKey, encryptionKey),
	}
}

func (store *CookieStore[T]) Load(ctx context.Context, session *Session) (*T, error) {
	result := new(T)

	if session.ID() == "" {
		return result, nil
	}

	m := map[any]any{}
	if err := securecookie.DecodeMulti(SessionIDKey, session.ID(), &m, store.codecs...); err != nil {
		return nil, err
	}

	if err := xreflect.Parse(result, false, "session", func(field string) any { return m[field] }); err != nil {
		return nil, err
	}

	return result, nil
}

func (store *CookieStore[T]) Save(ctx context.Context, session *Session, t *T) error {
	m := xreflect.ToMap(t, "session")
	encoded, err := securecookie.EncodeMulti(SessionIDKey, m, store.codecs...)
	if err != nil {
		return err
	}

	session.SetID(encoded)
	return nil
}
