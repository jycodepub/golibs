package dstruct

type Stack[T any] struct {
	stack []T
	size  int
	len   int
}
type StackFull struct{}

func (s *StackFull) Error() string {
	return "stack is full"
}

type StackEmpty struct{}

func (s *StackEmpty) Error() string {
	return "stack is empty"
}

func NewStack[T any](size int) *Stack[T] {
	return &Stack[T]{
		stack: make([]T, 0, size),
		size:  size,
		len:   0,
	}
}

func (s *Stack[T]) Push(e T) *StackFull {
	if s.len == s.size {
		return &StackFull{}
	}
	s.stack = append(s.stack, e)
	s.len += 1
	return nil
}

func (s *Stack[T]) Pop() (T, *StackEmpty) {
	if s.len == 0 {
		var t T
		return t, &StackEmpty{}
	}
	s.len -= 1
	return s.stack[s.len], nil
}

func (s *Stack[T]) IsEmpty() bool {
	return s.len == 0
}
