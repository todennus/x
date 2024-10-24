package lock

import (
	"context"
	"errors"
)

var (
	ErrLock = errors.New("lock_fails")
)

type Locker interface {
	Lock(context.Context) error
	Unlock(context.Context)
}

func Func(locker Locker, ctx context.Context, f func() error) error {
	if err := locker.Lock(ctx); err != nil {
		return err
	}

	defer locker.Unlock(ctx)

	return f()
}
