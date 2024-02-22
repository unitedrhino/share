package utils

import (
	"github.com/spf13/cast"
	"strings"
)

func GetIDPath(idPath string) (ret []int64) {
	ids := strings.Split(idPath, "-")
	for _, v := range ids {
		if v != "" {
			ret = append(ret, cast.ToInt64(v))
		}
	}
	return ret
}
func GetNamePath(namePath string) (ret []string) {
	ids := strings.Split(namePath, "-")
	for _, v := range ids {
		if v != "" {
			ret = append(ret, v)
		}
	}
	return ret
}
