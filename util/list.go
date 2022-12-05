package util

func ListToString(lst []string) string {
	str := "{"
	for i, e := range lst {
		str += e
		if i < len(lst)-1 {
			str += ","
		}
	}
	str = str + "}"
	return str
}

func Filter[T any](p func(T) bool, lst *[]T) {
	// filtered := (*lst)[:0]
	// for _, x := range *lst {
	// 	if p(x) {
	// 		filtered = append(filtered, x)
	// 	}
	// }
	// *lst = filtered
	i := 0
	for _, x := range *lst {
		if p(x) {
			(*lst)[i] = x
			i++
		}
	}
	*lst = (*lst)[:i]
}

func Filter2[T any](p func(T) bool, lst []T) []T {
	result := make([]T, 0, 3)
	for _, x := range lst {
		if p(x) {
			result = append(result, x)
		}
	}
	return result
}

func Contains[T comparable](target T, lst []T) bool {
	for _, x := range lst {
		if x == target {
			return true
		}
	}
	return false
}
