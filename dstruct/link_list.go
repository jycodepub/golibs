package dstruct

type Node[T any] struct {
	Value T
	Next *Node[T]
}

func NewNode[T any] (value T) *Node[T] {
	return &Node[T]{
		Value: value,
	}
}

type LinkList[T any] struct {
	Head *Node[T]
	Tail *Node[T]
}

func (l *LinkList[T]) Append(n *Node[T]) {
	if l.Head == nil && l.Tail == nil {
		l.Head = n
		l.Tail = n
	} else {
		l.Tail.Next = n
		l.Tail = n
	}
}

func (l *LinkList[T]) PopHead() *Node[T] {
	if l.Head == nil && l.Tail == nil {
		return nil
	}
	if l.Head != l.Tail {
		h := l.Head
		l.Head = h.Next
		return h
	} else {
		h := l.Head
		l.Head = nil
		l.Tail = nil	
		return h
	}
}

func (l *LinkList[T]) IsEmpty() bool {
	return l.Head == nil && l.Tail == nil
}