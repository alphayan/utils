package service

import "strconv"

func ArrayToSting(arrayData []int64) string {
	str := "("
	for k, v := range arrayData {
		if k == len(arrayData)-1 {
			str += strconv.FormatInt(v, 10) + ")"
		} else {
			str += strconv.FormatInt(v, 10) + ","
		}
	}
	return str
}

//uint32转int64
func Uint32ToInt64(source uint32) int64 {
	str := strconv.FormatUint(uint64(source), 10)
	des, _ := strconv.ParseInt(str, 10, 64)
	return des
}

//int到float32
func IntToFloat32(source int) float32 {
	return float32(source)
}

//字符串转int64
func StringToInt64(source string) int64 {
	des, _ := strconv.ParseInt(source, 10, 64)
	return des
}

//字符串转float64
func StringToFloat(source string) float64 {
	des, _ := strconv.ParseFloat(source, 54)
	return des
}

//int64转Int
func Int64ToInt(source int64) int {
	str := strconv.FormatInt(source, 10)
	des, _ := strconv.Atoi(str)
	return des
}