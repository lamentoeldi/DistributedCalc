package redis

import (
	"context"
	"github.com/distributed-calc/v1/pkg/redis"
	"testing"
	"time"
)

func TestBlacklist_Add(t *testing.T) {
	cases := []struct {
		name    string
		tokenID string
		ttl     time.Duration
		wantErr bool
	}{
		{
			name:    "success",
			tokenID: "1",
			ttl:     10 * time.Second,
			wantErr: false,
		},
	}

	client, err := redis.NewRedis(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	bl := NewBlacklist(client)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			err := bl.Add(ctx, tc.tokenID, tc.ttl)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error, got none")
			}
		})
	}
}

func TestBlacklist_Remove(t *testing.T) {
	cases := []struct {
		name    string
		tokenID string
		wantErr bool
	}{
		{
			name:    "success",
			tokenID: "2",
			wantErr: false,
		},
	}

	client, err := redis.NewRedis(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	bl := NewBlacklist(client)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err = bl.Add(ctx, "2", time.Minute)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			err := bl.Remove(ctx, tc.tokenID)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error, got none")
			}
		})
	}
}

func TestBlacklist_IsBlackListed(t *testing.T) {
	cases := []struct {
		name     string
		tokenID  string
		expected bool
		wantErr  bool
	}{
		{
			name:     "blacklisted",
			tokenID:  "3",
			expected: true,
		},
		{
			name:    "not-blacklisted",
			tokenID: "4",
		},
	}

	client, err := redis.NewRedis(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	bl := NewBlacklist(client)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	err = bl.Add(ctx, "3", time.Minute)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
			defer cancel()

			isBlacklisted, err := bl.IsBlackListed(ctx, tc.tokenID)
			if tc.wantErr == false && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.wantErr == true && err == nil {
				t.Errorf("expected error, got none")
			}

			if isBlacklisted != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, isBlacklisted)
			}
		})
	}
}
