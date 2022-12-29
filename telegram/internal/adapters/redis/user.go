package redis

import (
	"context"
	"encoding/json"

	"github.com/escalopa/govisa/telegram/core"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
)

type UserCache struct {
	r *redis.Client
}

// NewUserCache creates a new user cache instance with a redis client
func NewUserCache(rc *redis.Client) (*UserCache, error) {
	if rc == nil {
		return nil, errors.New("redis client is nil")
	}
	return &UserCache{r: rc}, nil
}

// GetUserByID Get user from cache by chat id
func (uc *UserCache) GetUserByID(ctx context.Context, id string) (*core.User, error) {
	data := uc.r.Get(ctx, id)
	if data.Err() != nil {
		return nil, data.Err()
	}
	user := &core.User{}
	if err := json.Unmarshal([]byte(data.Val()), user); err != nil {
		return nil, err
	}
	return user, nil
}

// SaveUserByID Save user in cache by chat id
func (uc *UserCache) SaveUserByID(ctx context.Context, user *core.User) error {
	err := uc.r.Set(context.Background(), user.ID, user, 0)
	if err.Err() != nil {
		return errors.Wrap(err.Err(), "error saving user in db")
	}
	return nil
}
