package list

// List is a container struct.
type List[T interface{ Equal(T) bool }] struct {
	data []T
}

// NewList constructs a List.
func NewList[T interface{ Equal(T) bool }]() List[T] {
	return List[T]{
		data: make([]T, 0),
	}
}

// Contains checks whether given element is contained or not.
func (list List[T]) Contains(x T) bool {
	for _, v := range list.data {
		if v.Equal(x) {
			return true
		}
	}

	return false
}

// Add adds a element in list.
func (list *List[T]) Add(x T) {
	list.data = append(list.data, x)
}

// Remove removes element from the list.
func (list *List[T]) Remove(x T) {
	found := false
	for i, v := range list.data {
		if v.Equal(x) {
			list.data[0], list.data[i] = list.data[i], list.data[0]
			found = true

			break
		}
	}

	if found {
		list.data = list.data[1:]
	}
}

// Data returns list's data.
func (list List[T]) Data() []T {
	return list.data
}

// Length returns length of the List.
func (list List[T]) Length() int {
	return len(list.data)
}
