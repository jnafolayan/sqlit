package lib

type Iterator[T interface{}] struct {
	list   []T
	cursor int
}

func NewIterator[T interface{}](list []T) *Iterator[T] {
	return &Iterator[T]{list: list, cursor: -1}
}

func (i *Iterator[T]) Next() (v T) {
	i.cursor++
	if i.cursor < len(i.list) {
		v = i.list[i.cursor]
	}

	return
}

func (i *Iterator[T]) Peek() (v T) {
	if i.cursor < len(i.list)-1 {
		v = i.list[i.cursor+1]
	}

	return
}
