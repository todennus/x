package session

import "context"

type Store[T any] interface {
	Load(ctx context.Context, id *Session) (*T, error)
	Save(ctx context.Context, id *Session, t *T) error
}
