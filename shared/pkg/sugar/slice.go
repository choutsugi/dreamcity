package sugar

func Contains[K comparable](eles []K, value K) bool {
	for _, ele := range eles {
		if ele == value {
			return true
		}
	}
	return false
}

func Delete[K comparable](eles []K, value K) []K {
	var res []K
	for _, ele := range eles {
		if value != ele {
			res = append(res, ele)
		}
	}
	return res
}
