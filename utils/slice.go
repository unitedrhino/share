package utils

// SliceIn 判断 in 是否在 cmp 中
func SliceIn[T comparable](in T, cmp ...T) bool {
	for _, v := range cmp {
		if in == v {
			return true
		}
	}
	return false
}

// SliceIndex 获取 T[index] 的值，否则返回默认 defaul
func SliceIndex[T any](slice []T, index int, defaul T) T {
	if index >= 0 && index < len(slice) {
		return slice[index]
	} else {
		return defaul
	}
}

func ToAnySlice[t any](in []t) (ret []any) {
	for _, v := range in {
		ret = append(ret, v)
	}
	return
}

func GetAddSlice[t comparable](oldV, newV []t) []t {
	var oldMap = SliceToSet(oldV)
	var ret []t
	for _, v := range newV {
		if _, ok := oldMap[v]; !ok {
			ret = append(ret, v)
		}
	}
	return ret
}

func SliceToSet[t comparable](in []t) map[t]struct{} {
	var retM = map[t]struct{}{}
	for _, v := range in {
		retM[v] = struct{}{}
	}
	return retM
}

func ToSliceWithFunc[inT any, retT any](in []*inT, f func(in *inT) retT) (ret []retT) {
	if in == nil {
		return nil
	}
	for _, v := range in {
		ret = append(ret, f(v))
	}
	return ret
}

func AnyToSlice[t any](in []any) (ret []t) {
	for _, v := range in {
		ret = append(ret, v.(t))
	}
	return
}

func NewFillSlice[T any](num int, val T) []T {
	sli := make([]T, num)
	for i := range sli {
		sli[i] = val
	}
	return sli
}

func SliceDelete[T comparable](base []T, val T) []T {
	for i, v := range base {
		if v == val {
			return append(base[:i], base[i+1:]...)
		}
	}
	return base
}
func SliceReversal[T comparable](base []T) (ret []T) {
	for i := len(base) - 1; i >= 0; i-- {
		ret = append(ret, base[i])
	}
	return ret
}
