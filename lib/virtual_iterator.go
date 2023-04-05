package lib

type VirtualIterator[T interface{}] struct {
	cursor int
	length int
	getter func(int) T
}

func NewVirtualIterator[T interface{}](length int, getter func(int) T) *VirtualIterator[T] {
	return &VirtualIterator[T]{
		length: length,
		getter: getter,
		cursor: -1,
	}
}

func (vi *VirtualIterator[T]) Next() (v T) {
	vi.cursor++
	if vi.cursor < vi.length {
		v = vi.getter(vi.cursor)
	}

	return
}

func (vi *VirtualIterator[T]) Peek() (v T) {
	if vi.cursor < vi.length-1 {
		v = vi.getter(vi.cursor + 1)
	}

	return
}
