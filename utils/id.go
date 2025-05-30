package utils

import (
	"gitee.com/unitedrhino/share/def"
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

func GenIDPath(in ...int64) string {
	return strings.Join(cast.ToStringSlice(in), "-") + "-"
}

func GenSliceStr[a any](s []a) string {
	if len(s) == 0 {
		return ""
	}
	return "," + strings.Join(cast.ToStringSlice(s), ",") + ","
}

func StrGenStrSlice(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	strs := strings.Split(s, ",")
	var ret []string
	for _, str := range strs {
		if str == "" {
			continue
		}
		ret = append(ret, str)
	}
	return ret
}

func StrGenInt64Slice(s string) []int64 {
	if len(s) == 0 {
		return []int64{}
	}
	strs := strings.Split(s, ",")
	var ret []int64
	for _, str := range strs {
		if str == "" {
			continue
		}
		ret = append(ret, cast.ToInt64(str))
	}
	return ret
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
