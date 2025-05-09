package dstruct

const (
	OutOfBound = 0
	EmptyList  = 1
)

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type LinkList[T any] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

type LinkListError struct {
	errCode int
}

func (e LinkListError) Error() string {
	if e.errCode == OutOfBound {
		return "LinkListError: out of bound"
	} else if e.errCode == EmptyList {
		return "LinkListError: empty list"
	} else {
		return "LinkListError: unknown error"
	}
}

func (l *LinkList[T]) Append(n *Node[T]) {
	if l.IsEmpty() {
		l.head = n
		l.tail = n
	} else {
		l.tail.Next = n
		l.tail = n
	}
	l.length += 1
}

func (l *LinkList[T]) Pop() (*Node[T], error) {
	if l.length == 0 {
		return nil, LinkListError{
			errCode: EmptyList,
		}
	}

	var n *Node[T]
	if l.length == 1 {
		n = l.head
		l.head = nil
		l.tail = nil
	} else {
		n = l.tail
		np := l.advance(l.length - 1)
		np.Next = nil
		l.tail = np
	}
	l.length -= 1

	return n, nil
}

func (l *LinkList[T]) RemoveHead() (*Node[T], error) {
	if l.length == 0 {
		return nil, LinkListError{
			errCode: EmptyList,
		}
	}

	var n *Node[T]
	if l.length > 1 {
		h := l.head
		l.head = h.Next
		n = h
	} else {
		h := l.head
		l.head = nil
		l.tail = nil
		n = h
	}
	l.length -= 1

	return n, nil
}

func (l *LinkList[T]) Get(i int) (*Node[T], error) {
	if i < 0 || i >= l.Len() {
		return nil, LinkListError{
			errCode: OutOfBound,
		}
	}

	if l.length == 0 {
		return nil, LinkListError{
			errCode: EmptyList,
		}
	} else {
		return l.advance(i), nil
	}
}

func (l *LinkList[T]) IsEmpty() bool {
	return l.length == 0
}

func (l *LinkList[T]) Len() int {
	return l.length
}

func (l *LinkList[T]) advance(stp int) *Node[T] {
	n := l.head
	for _ = range stp {
		n = n.Next
	}
	return n
}
