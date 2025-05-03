package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Blacklist struct {
	client redis.UniversalClient
}

func NewBlacklist(client redis.UniversalClient) *Blacklist {
	return &Blacklist{client: client}
}

func (bl *Blacklist) Add(ctx context.Context, tokenID string, ttl time.Duration) error {
	err := bl.client.Set(ctx, tokenID, "", ttl).Err()

	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}

	return nil
}

func (bl *Blacklist) Remove(ctx context.Context, tokenID string) error {
	err := bl.client.Del(ctx, tokenID).Err()

	if err != nil {
		return fmt.Errorf("failed to remove token from blacklist: %w", err)
	}

	return nil
}

func (bl *Blacklist) IsBlackListed(ctx context.Context, tokenID string) (bool, error) {
	exists, err := bl.client.Exists(ctx, tokenID).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if token is blacklisted: %w", err)
	}

	return exists > 0, nil
}
