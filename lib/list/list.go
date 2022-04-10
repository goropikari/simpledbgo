package list

type List[T interface{ Equal(T) bool }] struct {
	data []T
}

func NewList[T interface{ Equal(T) bool }]() List[T] {
	return List[T]{
		data: make([]T, 0),
	}
}

func (list List[T]) Contains(x T) bool {
	for _, v := range list.data {
		if v.Equal(x) {
			return true
		}
	}

	return false
}

func (list *List[T]) Add(x T) {
	list.data = append(list.data, x)
}

func (list *List[T]) Remove(x T) {
	for i, v := range list.data {
		if v.Equal(x) {
			list.data[0], list.data[i] = list.data[i], list.data[0]

			break
		}
	}

	list.data = list.data[1:]
}
