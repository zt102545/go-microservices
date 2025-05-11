package utils

import (
	"sync"
	"time"
)

type ExpiringMutex[T comparable] struct {
	mutexes sync.Map
	ttl     time.Duration
}

type mutexEntry struct {
	mutex      *sync.Mutex
	expiration time.Time
}

func NewExpiringMutex[T comparable](ttl time.Duration) *ExpiringMutex[T] {
	em := &ExpiringMutex[T]{ttl: ttl}
	go em.cleanupExpiredEntries()
	return em
}

func (em *ExpiringMutex[T]) Lock(key T) {
	v, _ := em.mutexes.LoadOrStore(key, &mutexEntry{
		mutex:      &sync.Mutex{},
		expiration: time.Now().Add(em.ttl),
	})
	entry := v.(*mutexEntry)

	// 获取锁
	entry.mutex.Lock()
	entry.expiration = time.Now().Add(em.ttl)
}

func (em *ExpiringMutex[K]) Unlock(key K) {
	v, ok := em.mutexes.Load(key)
	if !ok {
		panic("unlock of unlocked mutex")
	}
	entry := v.(*mutexEntry)
	// 释放锁
	entry.mutex.Unlock()

	// 延长过期时间, 确保 Mutex 在被解锁后不会立即过期
	entry.expiration = time.Now().Add(em.ttl)
}

func (em *ExpiringMutex[T]) cleanupExpiredEntries() {
	// 定期清理过期项
	for {
		time.Sleep(5 * time.Second)
		now := time.Now()

		em.mutexes.Range(func(key, value interface{}) bool {
			entry := value.(*mutexEntry)
			if entry.expiration.Before(now) {
				em.mutexes.Delete(key)
			}
			return true
		})
	}
}
