package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	lockPrefix        = "pc:lock:"
	onlineUsersPrefix = "pc:online:"
	draftPrefix       = "pc:draft:"
)

type SyncCache struct {
	rdb *redis.Client
}

func NewSyncCache(rdb *redis.Client) *SyncCache {
	return &SyncCache{rdb: rdb}
}

func (s *SyncCache) LockPointCloud(ctx context.Context, pointCloudID, userID string, ttl time.Duration) (bool, error) {
	key := lockPrefix + pointCloudID
	ok, err := s.rdb.SetNX(ctx, key, userID, ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *SyncCache) UnlockPointCloud(ctx context.Context, pointCloudID, userID string) error {
	key := lockPrefix + pointCloudID
	val, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	if val != userID {
		return fmt.Errorf("lock held by another user")
	}
	return s.rdb.Del(ctx, key).Err()
}

func (s *SyncCache) GetOnlineUsers(ctx context.Context, pointCloudID string) ([]string, error) {
	key := onlineUsersPrefix + pointCloudID
	return s.rdb.SMembers(ctx, key).Result()
}

func (s *SyncCache) AddOnlineUser(ctx context.Context, pointCloudID, userID string) error {
	key := onlineUsersPrefix + pointCloudID
	return s.rdb.SAdd(ctx, key, userID).Err()
}

func (s *SyncCache) RemoveOnlineUser(ctx context.Context, pointCloudID, userID string) error {
	key := onlineUsersPrefix + pointCloudID
	return s.rdb.SRem(ctx, key, userID).Err()
}

func (s *SyncCache) SetDraftAnnotation(ctx context.Context, pointCloudID, userID, data string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s:%s", draftPrefix, pointCloudID, userID)
	return s.rdb.Set(ctx, key, data, ttl).Err()
}

func (s *SyncCache) GetDraftAnnotation(ctx context.Context, pointCloudID, userID string) (string, error) {
	key := fmt.Sprintf("%s%s:%s", draftPrefix, pointCloudID, userID)
	return s.rdb.Get(ctx, key).Result()
}
