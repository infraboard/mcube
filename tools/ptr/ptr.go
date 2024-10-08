package ptr

func GetValue[T any](s *T) T {
	if s == nil {
		var zero T
		return zero
	}

	return *s
}

func GetArrayValue[T any](ptrItems []*T) []T {
	items := []T{}
	for i := range ptrItems {
		item := ptrItems[i]
		if item != nil {
			items = append(items, *item)
		}
	}

	return items
}
