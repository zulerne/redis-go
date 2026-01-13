package resp

func ToStringSlice(arr []Value) []string {
	res := make([]string, 0, len(arr))

	for _, v := range arr {
		res = append(res, v.Value)
	}

	return res
}
