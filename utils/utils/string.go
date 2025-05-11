package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// 数据库表名去除符号
func SetTable(t string) string {
	t = strings.Replace(t, "`", "", -1)
	return strings.Replace(t, "\"", "", -1)
}

// 版本号比较；-1=version1<version2；1=version1>version2；0=相同
func CompareVersion(version1 string, version2 string) int {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")

	for i := 0; i < len(v1) || i < len(v2); i++ {
		num1, num2 := 0, 0
		if i < len(v1) {
			num1, _ = strconv.Atoi(v1[i])
		}
		if i < len(v2) {
			num2, _ = strconv.Atoi(v2[i])
		}

		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}
	return 0
}

// 临时 国际化函数
func LocaleString(en string, ru string, ratio string) string {
	if ratio == "ru" {
		return ru
	}

	return en
}

// 数组去重
func RemoveDuplicatesString(nums []string) []string {
	encountered := map[string]struct{}{}
	result := make([]string, 0, len(nums))

	for _, v := range nums {
		if _, ok := encountered[v]; !ok {
			encountered[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// 移除字符串数组中包含指定子串的元素
func RemoveStringsContaining(input []string, substring string) []string {
	var result []string
	for _, str := range input {
		if !strings.Contains(str, substring) {
			result = append(result, str)
		}
	}
	return result
}

func StrToInt64(val string) int64 {
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		fmt.Println("转换错误:", err)
		return 0
	}
	return num
}
