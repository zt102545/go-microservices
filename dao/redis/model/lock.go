package model

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go-microservices/dao/redis/key"
	"time"
)

type (
	LockRedis interface {
		// Lock 加锁
		Lock(ctx context.Context, lockKey string, timeout int64) error
		// UnLock 解锁
		UnLock(ctx context.Context, lockKey string) error
		// IsLock 判断是否被锁
		IsLock(ctx context.Context, lockKey string) (bool, error)
	}

	customLockRedis struct {
		rc *redis.Client
	}
)

func NewLockRedis(rc *redis.Client) LockRedis {
	return &customLockRedis{
		rc: rc,
	}
}

// Lock 加锁
func (l *customLockRedis) Lock(ctx context.Context, lockKey string, timeout int64) error {
	redisKey := fmt.Sprintf("%s%s", key.RedisLock, lockKey)
	ok, err := l.rc.SetNX(ctx, redisKey, "1", time.Duration(timeout)*time.Second).Result()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("[%s] lock failed", redisKey)
	}
	return nil
}

// UnLock 解锁
func (l *customLockRedis) UnLock(ctx context.Context, lockKey string) error {
	redisKey := fmt.Sprintf("%s%s", key.RedisLock, lockKey)
	_, err := l.rc.Del(ctx, redisKey).Result()
	if err != nil {
		return err
	}
	return nil
}

// IsLock 判断是否被锁
func (l *customLockRedis) IsLock(ctx context.Context, lockKey string) (bool, error) {
	redisKey := fmt.Sprintf("%s%s", key.RedisLock, lockKey)
	r, err := l.rc.Exists(ctx, redisKey).Result()
	if err != nil {
		return false, err
	}
	if r > 0 {
		return true, nil
	}
	return false, nil
}
