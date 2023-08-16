package utils

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

var (
	rmbRatio  int64 = 100
	dRmbRatio       = decimal.NewFromInt(rmbRatio)
)

// StringToYuan 字符串类型的人民币元单位值转换为元单位的 float 值
func StringToYuan(str string) float64 {
	dec, _ := decimal.NewFromString(str)
	yuan, _ := dec.Float64()
	return yuan
}

// StringCentToYuan 字符串类型的人民币分单位值转换为元单位的值
func StringCentToYuan(cent string) float64 {
	return CentToYuan(cast.ToInt64(cent))
}

// IntCentToYuan int 类型的人民币分单位值转换为元单位的值
func IntCentToYuan(cent int) float64 {
	return CentToYuan(int64(cent))
}

// CentToYuan int64 类型的人民币分单位值转换为元单位的值
func CentToYuan(cent int64) float64 {
	yuan, _ := decimal.NewFromInt(cent).Div(dRmbRatio).Float64()
	return yuan
}

// StringYuanToCent 字符串类型的人民币元单位的值转换为分单位的值
func StringYuanToCent(yuan string) int64 {
	return YuanToCent(cast.ToFloat64(yuan))
}

// YuanToCent 浮点类型的人民币元单位的值转换为分单位的值
func YuanToCent(yuan float64) int64 {
	return decimal.NewFromFloat(yuan).Mul(dRmbRatio).IntPart()
}
