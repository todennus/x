package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/todennus/x/xcontext"
)

var _ Locker = (*RedisLock)(nil)

type RedisLock struct {
	key        string
	client     *redis.Client
	exipration time.Duration
}

func NewRedisLock(redis *redis.Client, key string, exipration time.Duration) *RedisLock {
	return &RedisLock{
		client:     redis,
		key:        key,
		exipration: exipration,
	}
}

func (l *RedisLock) Lock(ctx context.Context) error {
	for {
		ok, err := l.client.SetNX(ctx, l.key, "", l.exipration).Result()
		if err != nil {
			return fmt.Errorf("%w: %s", ErrLock, err.Error())
		}

		if ok {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (l *RedisLock) Unlock(ctx context.Context) {
	_, err := l.client.Del(ctx, l.key).Result()
	if err != nil {
		xcontext.Logger(ctx).Warn("failed-to-release-redis-lock", "key", l.key, "err", err)
	}
}
