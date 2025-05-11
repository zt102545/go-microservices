package gos

import "github.com/panjf2000/ants/v2"

var GoPool *ants.Pool

const DefaultGoroutineCount = 1000

func InitGoroutinePool(num int) {
	var err error
	if num <= 0 {
		num = DefaultGoroutineCount
	}
	GoPool, err = ants.NewPool(num)
	if err != nil {
		panic(err)
	}
	return
}

func GetGoroutinePool(num int) *ants.Pool {
	var err error
	if num <= 0 {
		num = DefaultGoroutineCount
	}
	pool, err := ants.NewPool(num)
	if err != nil {
		panic(err)
	}
	return pool
}

func Tune(count int) {
	GoPool.Tune(count)
}

func Release() {
	GoPool.Release()
}
