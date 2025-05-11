package utils

import (
	"math/rand"
	"time"
)

// 根据权重值随机获取下标
func WeightedRandom(weights []int) int {

	total := 0
	for _, value := range weights {
		total += value
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := rand.Intn(total)
	for key, value := range weights {
		randNum -= value
		if randNum < 0 {
			return key
		}
	}
	return -1
}

// RandInt64 生成一个在[min, max]范围内的随机int64整数。
// 如果min大于等于max，则返回max。
// 参数:
//
//	min - 随机数生成的下限。
//	max - 随机数生成的上限。
//
// 返回值:
//
//	生成的随机整数，该整数在[min, max]范围内。
func RandInt64(min, max int64) int64 {
	// 检查min和max的值，以确保它们构成有效的范围。
	if min >= max {
		return max
	}
	rand.New(rand.NewSource(time.Now().UnixNano())) // 设置随机数种子
	// 使用rand包的Int63n函数生成[min, max]范围内的随机数。
	// 注意：Int63n返回的随机数是[0, max-min)范围内的，因此需要加上min以将其调整到正确的范围。
	return rand.Int63n(max-min) + min
}

// GenerateRandomInt 生成随机数，不包含max
func GenerateRandomInt(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano())) // 设置随机数种子
	// 计算范围内随机数
	return rand.Intn(max-min) + min
}
