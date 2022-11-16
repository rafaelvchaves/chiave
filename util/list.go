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
	i := 0
	for _, x := range *lst {
		if p(x) {
			(*lst)[i] = x
			i++
		}
	}
	*lst = (*lst)[:i]
}

func Contains[T comparable](target T, lst []T) bool {
	for _, x := range lst {
		if x == target {
			return true
		}
	}
	return false
}
