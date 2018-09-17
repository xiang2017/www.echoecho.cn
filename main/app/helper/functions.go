package helper

import (
	"strings"
)

func InStringArray(s string, arr *[]string, ignore bool) bool {
	if ignore {
		s = strings.ToLower(s)
	}
	for i := 0; i < len(*arr); i ++{
		if ignore {
			if s == strings.ToLower((*arr)[i]) {
				return true
			}
		} else{
			if s == (*arr)[i] {
				return true
			}
		}
	}
	return false
}

func RemoveStringFromArray (s string, arr *[]string, ignore bool) [] string{
	if ignore {
		s = strings.ToLower(s)
	}
	for i := 0; i < len(*arr); i ++{
		if ignore {
			if s == strings.ToLower((*arr)[i]) {
				return append((*arr)[0:i], (*arr)[i+1: len(*arr)]...)
			}
		} else{
			if s == (*arr)[i] {
				return append((*arr)[0:i], (*arr)[i+1: len(*arr)]...)
			}
		}
	}
	return *arr
}

func InIntArray(num int, arr *[]int) bool {
	for i := 0; i < len(*arr); i ++ {
		if (*arr)[i] == num {
			return true
		}
	}
	return false
}

func RemoveIndexFromArray (index int, arr *[]int) []int {
	return append((*arr)[0:index], (*arr)[index+1 : len(*arr)]...)
}