package utils

import (
	"fmt"
	"github.com/spf13/cast"
	"golang.org/x/exp/constraints"
	"math"
	"strconv"
	"strings"
)

type CanAdd interface {
	constraints.Integer | constraints.Float
}

func Sum[addT CanAdd](datas ...addT) (sum addT) {
	for _, v := range datas {
		sum += v
	}
	return
}

// 保留n位小数
func Decimal[valueType constraints.Float](value valueType, n int) valueType {
	if math.IsNaN(float64(value)) {
		return 0
	}
	v, _ := strconv.ParseFloat(fmt.Sprintf("%."+cast.ToString(n)+"f", value), 64)
	return valueType(v)
}
func Max[t CanAdd](in []t) t {
	if len(in) == 0 {
		return 0
	}
	var max t
	for _, v := range in {
		if v > max {
			max = v
		}
	}
	return max
}

func Min[t CanAdd](in []t) t {
	if len(in) == 0 {
		return 0
	}
	var min t
	for _, v := range in {
		if v < min {
			min = v
		}
	}
	return min
}

// 动态确定step的小数位数
func getDecimalPlaces(step float64) int {
	// 将float64转换为字符串，保留足够的精度
	str := strconv.FormatFloat(step, 'f', 16, 64)

	// 查找小数点位置
	if pos := strings.Index(str, "."); pos != -1 {
		// 计算小数点后的位数（去除末尾的0）
		decimalPart := str[pos+1:]
		// 从后往前找到第一个非零数字
		lastNonZero := len(decimalPart) - 1
		for lastNonZero >= 0 && decimalPart[lastNonZero] == '0' {
			lastNonZero--
		}
		return lastNonZero + 1
	}
	return 0 // 整数，没有小数部分
}

// 舍入到指定小数位数
func RoundToDecimal(num float64, decimal int) float64 {
	factor := math.Pow(10, float64(decimal))
	return math.Round(num*factor) / factor
}

// 处理数值的主函数
func StepFloat(num, step float64) float64 {
	if step == 0 || math.IsNaN(step) || math.IsInf(step, 0) {
		return num
	}

	// 执行Floor计算
	num = math.Floor(num/step) * step

	// 获取step的小数位数
	decimals := getDecimalPlaces(step)

	// 舍入到相应的精度
	return RoundToDecimal(num, decimals)
}
