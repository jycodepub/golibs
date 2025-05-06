package dstruct

type Node[T any] struct {
	Value T
	Next *Node[T]
}

type LinkList[T any] struct {
	head *Node[T]
	tail *Node[T]
	length int
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

func (l *LinkList[T]) PopHead() *Node[T] {
	if l.length == 0 {
		return nil
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
	return n
}

func (l * LinkList[T]) Get(i int) *Node[T] {
	if l.length == 0 {
		return nil
	} else {
		var n *Node[T] = l.head
		ri := 0
		for {
			if ri == i {
				break
			}
			n = n.Next
			ri += 1
		}
		return n
	}
}

func (l *LinkList[T]) IsEmpty() bool {
	return l.length == 0
}

func (l *LinkList[T]) Len() int {
	return l.length
}