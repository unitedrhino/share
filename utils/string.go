package utils

import (
	"strings"
	"unicode"
)

// 从目标串src中查找第n个目标字符c所在位置下标
func IndexN(src string, c byte, n int) int {
	var s []byte
	s = []byte(src)

	for i := 0; i < len(s); i++ {
		if n == 0 {
			return i
		}
		if s[i] == c {
			n--
		}
	}
	return -1
}

// SplitCutset 按数组 cuset 里的分隔符，对 str 进行切割
func SplitCutset(str, cutset string) []string {
	words := strings.FieldsFunc(str, func(r rune) bool {
		return strings.ContainsRune(cutset, r)
	})
	result := make([]string, 0, len(words))
	for _, w := range words {
		wd := strings.TrimSpace(w)
		if wd != "" {
			result = append(result, wd)
		}
	}
	return result
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func NewFillString(num int, val string, sep string) string {
	sli := NewFillSlice(num, val)
	return strings.Join(sli, sep)
}

// 单词全部转化为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// 单词全部转化为小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// 下划线单词转为大写驼峰单词
func UderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	if s != "Id" {
		s = strings.ReplaceAll(s, "Id", "ID")
	}
	return strings.Replace(s, " ", "", -1)
}

// 下划线单词转为小写驼峰单词
func UderscoreToLowerCamelCase(s string) string {
	s = UderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// 驼峰单词转下划线单词
func CamelCaseToUdnderscore(s string) string {
	var output []rune
	s = strings.ReplaceAll(s, "ID", "Id")
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
		} else {
			if unicode.IsUpper(r) {
				output = append(output, '_')
			}

			output = append(output, unicode.ToLower(r))
		}
	}
	return string(output)
}
