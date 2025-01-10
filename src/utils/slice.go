package utils

func Filter[E any](src []E, callback func(E) bool) []E {
	res := make([]E, 0, len(src))

	for i := 0; i < len(src); i++ {
		srci := src[i]
		if callback(srci) {
			res = append(res, srci)
		}
	}

	return res
}
