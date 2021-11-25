package util

func MaxInt64(num ...int64) int64 {
	var ret int64
	ret = num[0]
	for _, v := range num {
		if ret < v {
			ret = v
		}
	}
	return ret
}
