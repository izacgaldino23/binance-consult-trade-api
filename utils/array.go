package utils

type Slice[T any] []T

func (l *Slice[T]) Each(callback func(i int, v *T)) {
	for index := range *l {
		callback(index, &(*l)[index])
	}
}

func (l *Slice[T]) Filter(callback func(v *T, i int) bool) (res []T) {
	res = make([]T, 0)

	for index := range *l {
		if callback(&(*l)[index], index) {
			res = append(res, (*l)[index])
		}
	}

	return res
}

func (l *Slice[T]) Find(callback func(v *T, i int) bool) (res T) {
	for index := range *l {
		if callback(&(*l)[index], index) {
			res = (*l)[index]
			break
		}
	}

	return res
}
