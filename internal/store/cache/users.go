package cache

import (
	"SOCIAL/internal/store"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

// Gets the user from the cache
func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	fmt.Printf("Getting user from cache: %v", userID)
	cacheKey := fmt.Sprintf("user-%d", userID)
	fmt.Printf("Cache key: %v", cacheKey)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

// Captures the information in the cache
func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, UserExpTime).Err()

}
