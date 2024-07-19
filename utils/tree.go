package utils

import (
	"gitee.com/i-Things/share/def"
	"github.com/spf13/cast"
	"strings"
)

func IDPathHasAcess(idPath string, id int64) bool {
	if id == def.RootNode {
		return true
	}
	path := GetIDPath(idPath)
	if SliceIn(id, path...) {
		return true
	}
	return false
}

func GenIDPath(in []int64) string {
	return strings.Join(cast.ToStringSlice(in), "-") + "-"
}

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
