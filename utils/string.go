package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"gitee.com/unitedrhino/share/errors"
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

func JoinWithFunc[t any](elems []*t, sep string, f func(in *t) string) string {
	var strs []string
	for _, v := range elems {
		strs = append(strs, f(v))
	}
	return strings.Join(strs, sep)
}

// ParseNumberString 解析包含数字和数字范围的字符串，返回int64数组
// 支持的格式：[1,3,6] 或 [1-2,3-7] 或混合格式 [1,3-5,7]
func ParseNumberString(s string) ([]int64, error) {
	// 移除字符串前后的空白字符
	s = strings.TrimSpace(s)

	// 验证输入格式是否正确
	re := regexp.MustCompile(`^\[([\d,-]+)\]$`)
	if !re.MatchString(s) {
		return nil, errors.Parameter.AddMsg("invalid format, expected [num1,num2-num3,...]")
	}

	// 提取括号内的内容
	content := re.FindStringSubmatch(s)[1]
	if content == "" {
		return []int64{}, nil
	}

	// 按逗号分割内容
	parts := strings.Split(content, ",")
	var result []int64

	// 处理每个部分
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 检查是否是范围格式
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			// 解析范围的起始和结束数字
			start, err1 := strconv.ParseInt(strings.TrimSpace(rangeParts[0]), 10, 64)
			end, err2 := strconv.ParseInt(strings.TrimSpace(rangeParts[1]), 10, 64)

			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid number in range: %s", part)
			}

			if start > end {
				return nil, fmt.Errorf("start number greater than end in range: %s", part)
			}

			// 添加范围内的所有数字
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
		} else {
			// 解析单个数字
			num, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", part)
			}
			result = append(result, num)
		}
	}

	return result, nil
}
