package tools

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

func GenerateMd5(x int) string {
	md5String := fmt.Sprintf("%d", x)
	md5Byte := []byte(md5String)
	resByte := md5.Sum(md5Byte)
	resString := fmt.Sprintf("%x", resByte)
	return resString
}

func StringToInt(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		result = 0
	}
	return result
}

func Max(x int, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

func Min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}
