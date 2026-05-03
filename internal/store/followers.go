package store

import (
	"context"
	"database/sql"
	"log"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, followerID, userID)
	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, unfollowedUserId, followerUserID int64) error {
	log.Printf("Unfollow: %d, %d", followerUserID, unfollowedUserId)
	query := `
    DELETE from followers
    WHERE user_id = $1 AND follower_id = $2
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	log.Println(query)
	_, err := s.db.ExecContext(ctx, query, unfollowedUserId, followerUserID)
	return err
}
