package utils

import (
	"sort"
	"strconv"
)

// 是否存在数组中
func IsInArray[T comparable](target T, arr []T) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

// InsertAt 在切片中指定位置插入一个元素，并返回新切片。
// 参数 slice 是原始切片，index 是插入位置，element 是要插入的元素。
// 如果 index 超出切片的范围，将元素插入到切片末尾。
// 返回值是包含新元素的切片。
//
// 该函数不修改原切片，而是创建一个新的切片以包含插入的元素。
// 这种方式适用于需要保持原切片不变的场景，比如在数据结构中插入新元素。
func InsertAt[T any](slice []T, index int, element T) []T {
	// 检查 index 是否在有效范围内，如果不是，则将其设置为切片的长度。
	// 这样做是为了将元素插入到切片的末尾。
	if index < 0 || index > len(slice) {
		index = len(slice)
	}

	// 创建一个新的切片，长度比原切片多1，用于存放插入的元素。
	newSlice := make([]T, len(slice)+1)
	// 复制原切片的前 index 个元素到新切片。
	copy(newSlice, slice[:index])
	// 在新切片的 index 位置插入新元素。
	newSlice[index] = element
	// 复制原切片从 index 开始的元素到新切片的 index+1 位置开始。
	copy(newSlice[index+1:], slice[index:])
	// 返回包含新元素的新切片。
	return newSlice
}

// FindIndex 在切片中查找指定值的第一个出现位置。
// 如果找到该值，返回其索引；如果未找到，返回 -1。
// 参数:
//
//	slice []T: 要搜索的切片。
//	val T: 要查找的值。
//
// 返回值:
//
//	int: 查找值的索引，如果未找到则为 -1。
//
// 使用类型参数 T，使该函数适用于任何可比较的类型。
func FindIndex[T comparable](slice []T, val T) int {
	// 遍历切片中的每个元素。
	for i, v := range slice {
		// 如果当前元素与查找值相等，返回当前索引。
		if v == val {
			return i
		}
	}
	// 如果遍历结束都没有找到匹配的值，返回 -1 表示未找到。

	return -1 // 表示元素未找到
}

// RemoveElementByIndex 通过索引从切片中移除一个元素。
// 这个函数接受一个切片和一个索引作为参数，返回一个新的切片，其中不包含指定索引处的元素。
// 参数:
//
//	slice []T: 要操作的切片，可以是任何类型。
//	index int: 要移除元素的索引。
//
// 返回值:
//
//	[]T: 移除指定元素后的新切片。
//
// 说明:
//
//	这个函数不修改原始切片，而是通过创建一个新的切片来实现元素的移除，以保持函数的纯净性。
func RemoveElementByIndex[T any](slice []T, index int) []T {
	// 使用append函数和切片操作来构建一个新的切片，其中不包含指定索引处的元素。
	// 通过将原切片的前半部分和后半部分（不包括索引处的元素）连接起来实现移除操作。
	return append(slice[:index], slice[index+1:]...)
}

// RemoveElement 移除切片中所有与指定元素相同的元素，返回一个新的切片。
// 参数 slice 是原始的切片，element 是需要移除的元素。
// 返回值是一个新的切片，其中不包含指定的元素。
// 该函数通过遍历原始切片，并将不等于指定元素的元素追加到一个新的切片中，从而实现移除功能。
// 使用类型参数 T 使得该函数可以适用于任何可比较的类型。
func RemoveElement[T comparable](slice []T, element T) []T {
	// 初始化一个空的切片，用于存放移除指定元素后的结果。
	// 使用append函数和切片操作来构建一个新的切片，其中不包含指定索引处的元素。
	// 通过将原切片的前半部分和后半部分（不包括索引处的元素）连接起来实现移除操作。
	var result []T

	// 遍历原始切片，检查每个元素。
	for _, v := range slice {
		// 如果当前元素不等于指定元素，则将其追加到结果切片中。
		if v != element {
			result = append(result, v)
		}
	}

	// 返回移除指定元素后的新切片。
	return result
}

// 数组去重
func RemoveDuplicates[T comparable](nums []T) []T {
	encountered := map[T]struct{}{}
	result := make([]T, 0, len(nums))

	for _, v := range nums {
		if _, ok := encountered[v]; !ok {
			encountered[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// 两个数组取交集
func Intersection(nums1 []int64, nums2 []int64) []int64 {
	set := make(map[int64]struct{})
	result := make([]int64, 0)

	for _, num := range nums1 {
		set[num] = struct{}{}
	}

	for _, num := range nums2 {
		if _, ok := set[num]; ok {
			result = append(result, num)
			delete(set, num)
		}
	}

	return result
}

// 字符串数组排序,默认正序
func SortStringSlice(input []string, order ...string) []string {
	if len(order) > 0 && order[0] == "desc" {
		sort.Sort(sort.Reverse(sort.StringSlice(input)))
	} else {
		sort.Strings(input)
	}
	return input
}

// int64数组排序,默认正序
func SortInt64Slice(input []int64, order ...string) []int64 {
	if len(order) > 0 && order[0] == "desc" {
		sort.Slice(input, func(i, j int) bool {
			return input[i] > input[j]
		})
	} else {
		sort.Slice(input, func(i, j int) bool {
			return input[i] < input[j]
		})
	}
	return input
}

// AsMap 根据指定的回调函数将列表中的元素转换为键值对映射。
// 这个函数支持多种类型的键，包括int、int64、string、float32和float64。
// T2可以是任何类型，表示列表中元素的类型。
// 参数list是待处理的元素列表。
// callback是一个函数，它接受一个T2类型的指针，并返回一个T1类型的值，这个值将用作映射的键。
// 返回值是一个映射，其中键是通过回调函数从列表元素中获得的，值是指向原始列表元素的指针。
func AsMap[T1 int | int64 | string | float32 | float64 | uint64 | uint, T2 any](list []T2, callback func(*T2) T1) map[T1]*T2 {
	// 初始化一个空的映射，用于存储转换后的键值对。
	m := make(map[T1]*T2)

	// 如果列表为空，直接返回空映射，无需进一步处理。
	if len(list) == 0 {
		return m
	}

	// 遍历列表中的每个元素。
	for _, v := range list {
		// 使用回调函数将当前元素转换为键，并将对应的值（即元素的指针）存储到映射中。
		m[callback(&v)] = &v
	}

	// 返回构建完成的映射。
	return m
}

// Pluck 通过回调函数从列表中提取特定属性或值。
// list 是要处理的元素列表。
// callback 是一个函数，用于从每个元素中提取所需信息并返回。
// 返回一个包含通过回调函数从每个元素中提取的值的切片。
// T1 是回调函数返回值的类型，可以是int、int64、string、float32或float64之一。
// T2 是列表中元素的类型，可以是任何类型。
func Pluck[T1 int | int64 | string | float32 | float64 | uint64, T2 any](list []T2, callback func(*T2) T1) []T1 {
	// 初始化一个空的结果切片，用于存储回调函数的返回值。
	var result []T1

	// 如果列表为空，则直接返回空的结果切片。
	if len(list) == 0 {
		return result
	}

	// 遍历列表中的每个元素。
	for _, v := range list {
		// 将回调函数应用于当前元素，并将结果追加到结果切片中。
		result = append(result, callback(&v))
	}

	// 返回包含所有回调函数返回值的结果切片。
	return result
}

// 遍历 slice a，对于每个元素，检查它是否存在于 slice b 的 map 中。
// 如果不存在，则将其添加到结果切片中。
func Pluck2[T1 any, T2 any](list []T2, callback func(*T2) T1) []T1 {
	// 初始化一个空的结果切片，用于存储回调函数的返回值
	var result []T1

	// 如果列表为空，则直接返回空的结果切片
	if len(list) == 0 {
		return result
	}

	// 遍历列表中的每个元素
	for _, v := range list {
		// 将回调函数应用于当前元素，并将结果追加到结果切片中
		result = append(result, callback(&v))
	}

	// 返回包含所有回调函数返回值的结果切片
	return result
}

func ArrayDiff[T comparable](a, b []T) []T {
	// 创建一个映射用于存储b中的元素
	bMap := make(map[T]bool)
	for _, value := range b {
		bMap[value] = true
	}

	// 创建一个结果切片，用于存放差集
	var diff []T
	for _, value := range a {
		if !bMap[value] {
			diff = append(diff, value)
		}
	}

	return diff
}

// 子集检查函数函数，用于检查 subset 是否为 superset 的子集
func IsSubset[T comparable](subset, superset []T) bool {
	// 创建一个 map 来存储 superset 中的元素
	set := make(map[T]bool)
	for _, v := range superset {
		set[v] = true
	}

	// 遍历 subset，确保所有元素都在 superset 中
	for _, v := range subset {
		if !set[v] {
			return false // 只要有一个元素不在 superset 中，则返回 false
		}
	}

	return true
}

// FindFirst 在给定的列表中找到第一个满足回调函数条件的元素。
// 参数:
//
//	list []T: 待搜索的元素列表。
//	callback func(*T) bool: 回调函数，用于检查每个元素是否满足特定条件。
//
// 返回值:
//
//	*T: 如果找到满足条件的元素，返回该元素的指针；否则返回nil。
//
// 该函数通过遍历列表并应用回调函数来快速定位满足特定条件的第一个元素。
// 这种设计允许在不改变原始列表的情况下，灵活地对列表中的元素进行条件检查。
func FindFirst[T any](list []T, callback func(*T) bool) *T {
	// 遍历列表中的每个元素
	for _, v := range list {
		// 使用回调函数检查当前元素是否满足条件
		if callback(&v) {
			// 如果满足条件，返回当前元素的指针
			return &v
		}
	}

	// 如果没有找到满足条件的元素，返回nil
	return nil
}

// FindFirstPoint 在给定的切片中找到第一个满足回调函数条件的元素。
// 参数列表：
//   - list: 一个指向任意类型元素的指针切片。
//   - callback: 一个函数，用于检查列表中的每个元素是否满足某种条件。
//
// 返回值：
//   - 如果找到满足条件的元素，则返回该元素的指针；否则返回nil。
//
// 该函数遍历列表，对每个元素应用回调函数，一旦找到满足条件的元素，即立即返回该元素。
// 这种方法适用于在大型数据集中快速定位满足特定条件的元素，而无需检查所有元素。
func FindFirstPoint[T any](list []*T, callback func(*T) bool) *T {
	// 遍历列表中的每个元素
	for _, v := range list {
		// 应用回调函数检查当前元素是否满足条件
		if callback(v) {
			// 如果满足条件，立即返回当前元素的指针
			return v
		}
	}

	// 如果没有找到满足条件的元素，返回nil
	return nil
}

// GroupBy 对给定的切片进行分组，根据回调函数返回的键值。
// T: 切片中的元素类型。
// K: 回调函数返回的键值类型，必须是可比较的。
// list: 需要分组的切片。
// callback: 用于生成每个元素分组键值的回调函数。
// 返回一个映射，其中键是回调函数生成的键值，值是对应键值的元素切片。
func GroupBy[T any, K comparable](list []T, callback func(*T) K) map[K][]*T {
	// 初始化返回的映射。
	retMap := make(map[K][]*T)
	// 遍历切片中的每个元素。
	for _, v := range list {
		// 使用回调函数生成当前元素的键值。
		key := callback(&v)
		// 将当前元素添加到对应键值的切片中。
		// 如果键不存在，则创建一个新的切片。
		retMap[key] = append(retMap[key], &v)
	}
	// 返回分组后的映射。
	return retMap
}

// Filter根据回调函数过滤列表中的元素。
// 它接收一个类型为T的切片和一个函数，该函数接受一个T类型的指针并返回一个布尔值。
// Filter返回一个新的切片，其中包含原始切片中所有使回调函数返回true的元素的指针。
//
// 参数:
//
//	list []T: 待过滤的切片。
//	callback func(*T) bool: 用于决定元素是否应被包含在返回切片中的回调函数。
//
// 返回值:
//
//	[]*T: 过滤后的元素指针切片。
func Filter[T any](list []T, callback func(*T) bool) []*T {
	var retList []*T

	// 遍历输入切片中的每个元素。
	for _, v := range list {
		// 如果回调函数对当前元素返回true，则将该元素的指针添加到返回切片中。
		if callback(&v) {
			retList = append(retList, &v)
		}
	}

	// 返回过滤后的元素指针切片。
	return retList
}

// FilterPoint 根据回调函数过滤元素列表。
// 该函数接受一个指向T类型的指针切片和一个回调函数作为参数，回调函数用于判断每个元素是否应该被包含在返回的切片中。
// FilterPoint 返回一个新的切片，其中包含了原切片中所有使回调函数返回true的元素。
// T 可以是任何类型，这得益于Go的泛型特性，使得该函数可以适用于多种类型的列表。
//
// 参数:
//
//	list: 待过滤的元素指针切片。
//	callback: 用于决定元素是否被包含在返回切片中的回调函数。
//	   回调函数接受一个指向T类型的指针作为参数，并返回一个bool值。
//	   如果回调函数返回true，则表示对应的元素应该被包含在返回的切片中。
//
// 返回值:
//
//	包含原切片中满足回调函数条件的元素的指针切片。
func FilterPoint[T any](list []*T, callback func(*T) bool) []*T {
	// 初始化一个空的切片，用于存储过滤后的元素。
	var retList []*T

	// 遍历输入的列表。
	for _, v := range list {
		// 使用回调函数检查当前元素是否应该被包含在返回的切片中。
		if callback(v) {
			// 如果回调函数返回true，则将当前元素添加到返回的切片中。
			retList = append(retList, v)
		}
	}

	// 返回包含过滤后元素的切片。
	return retList
}

// ListMax 寻找列表中通过回调函数评估的最大值。
// 使用泛型允许适用于任何类型的列表。
// 参数 list 是一个泛型列表，包含了待评估的元素。
// 参数 callback 是一个函数指针，它接受列表中元素的指针，并返回一个int64类型的值，这个值用于比较和确定最大值。
// 返回值是列表中通过回调函数评估出的最大值。
func ListMax[T any](list []T, callback func(*T) int64) int64 {
	var maxNum int64 = 0

	// 遍历列表中的每个元素
	for _, v := range list {
		// 调用回调函数对当前元素进行评估，并将结果存储在ret中
		ret := callback(&v)

		// 如果当前元素的评估结果大于当前最大值，则更新最大值
		if ret > maxNum {
			maxNum = ret
		}
	}

	// 返回找到的最大值
	return maxNum
}

// Int64ToString int64数组转字符串数组
// Int64ToString 将一个int64类型的切片转换为对应字符串类型的切片。
// 这个函数的存在是因为int64类型的数据在某些情况下需要以字符串的形式进行处理或展示。
// 参数:
//
//	input []int64: 需要转换的int64类型切片。
//
// 返回值:
//
//	[]string: 转换后的字符串类型切片。
func Int64ToString(input []int64) []string {
	// 创建一个长度与input相同的字符串切片，用于存放转换后的结果。
	result := make([]string, len(input))

	// 遍历input切片，将每个int64类型的值转换为字符串类型，并存放到result切片中。
	for i, v := range input {
		// 使用strconv包的FormatInt函数将int64类型的值转换为字符串。
		// 参数1: v, 需要转换的int64类型值。
		// 参数2: 10, 指定转换为十进制形式的字符串。
		result[i] = strconv.FormatInt(v, 10)
	}

	// 返回转换后的字符串切片。
	return result
}

func Keys[T1 comparable, T2 any](m map[T1]T2) []T1 {
	var keys []T1
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[T1 comparable, T2 any](m map[T1]T2) []T2 {
	var values []T2
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
