package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func SnakeToCamel(s string) string {
	re := regexp.MustCompile("_(.)")
	return re.ReplaceAllStringFunc(s, func(m string) string {
		return strings.ToUpper(m[1:])
	})
}
func FirstToLower(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
func FirstLetterUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	// 将字符串转换为rune切片以便处理Unicode字符
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
func FirstLetterLower(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}
func UniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//	func LastNEquals(s string, substr string) bool {
//		n := len(substr)
//		if len(s) < n {
//			return false
//		}
//		return s[len(s)-n:] == substr
//	}
func SplitString(s string, n int) (string, string) {
	if n <= 0 {
		return "", s
	}
	if n >= len(s) {
		return s, ""
	}
	return s[:n], s[n:]
}
