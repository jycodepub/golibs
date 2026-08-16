package dstruct

import (
	"testing"
)

func TestLinkList(t *testing.T) {
	t.Run("New LinkList initial state", func(t *testing.T) {
		l := &LinkList[int]{}
		if !l.IsEmpty() {
			t.Errorf("expected new link list to be empty")
		}
		if l.Len() != 0 {
			t.Errorf("expected Len to be 0, got %d", l.Len())
		}
	})

	t.Run("Append and Len", func(t *testing.T) {
		l := &LinkList[string]{}
		n1 := &Node[string]{Value: "alpha"}
		n2 := &Node[string]{Value: "beta"}

		l.Append(n1)
		if l.Len() != 1 || l.IsEmpty() {
			t.Errorf("expected len 1, got %d", l.Len())
		}

		l.Append(n2)
		if l.Len() != 2 {
			t.Errorf("expected len 2, got %d", l.Len())
		}
	})

	t.Run("Get element", func(t *testing.T) {
		l := &LinkList[int]{}
		l.Append(&Node[int]{Value: 10})
		l.Append(&Node[int]{Value: 20})
		l.Append(&Node[int]{Value: 30})

		n0, err := l.Get(0)
		if err != nil || n0.Value != 10 {
			t.Errorf("expected 10 at index 0, got %v (err: %v)", n0, err)
		}

		n1, err := l.Get(1)
		if err != nil || n1.Value != 20 {
			t.Errorf("expected 20 at index 1, got %v (err: %v)", n1, err)
		}

		n2, err := l.Get(2)
		if err != nil || n2.Value != 30 {
			t.Errorf("expected 30 at index 2, got %v (err: %v)", n2, err)
		}

		// Out of bound checks
		_, err = l.Get(-1)
		if err == nil || err.Error() != "LinkListError: out of bound" {
			t.Errorf("expected out of bound error for -1 index, got %v", err)
		}

		_, err = l.Get(3)
		if err == nil || err.Error() != "LinkListError: out of bound" {
			t.Errorf("expected out of bound error for index 3, got %v", err)
		}
	})

	t.Run("Pop elements", func(t *testing.T) {
		l := &LinkList[int]{}
		l.Append(&Node[int]{Value: 100})
		l.Append(&Node[int]{Value: 200})

		// Pop last
		node, err := l.Pop()
		if err != nil || node.Value != 200 {
			t.Errorf("expected popped node 200, got %v (err: %v)", node, err)
		}
		if l.Len() != 1 {
			t.Errorf("expected length 1, got %d", l.Len())
		}

		// Pop remaining (single element)
		node, err = l.Pop()
		if err != nil || node.Value != 100 {
			t.Errorf("expected popped node 100, got %v (err: %v)", node, err)
		}
		if l.Len() != 0 || !l.IsEmpty() {
			t.Errorf("expected empty list after popping all")
		}

		// Pop from empty list
		_, err = l.Pop()
		if err == nil || err.Error() != "LinkListError: empty list" {
			t.Errorf("expected empty list error when popping empty list, got %v", err)
		}
	})

	t.Run("RemoveHead elements", func(t *testing.T) {
		l := &LinkList[string]{}
		l.Append(&Node[string]{Value: "first"})
		l.Append(&Node[string]{Value: "second"})

		// Remove head from multi-element list
		node, err := l.RemoveHead()
		if err != nil || node.Value != "first" {
			t.Errorf("expected head node 'first', got %v (err: %v)", node, err)
		}
		if l.Len() != 1 {
			t.Errorf("expected len 1, got %d", l.Len())
		}

		// Remove head from single-element list
		node, err = l.RemoveHead()
		if err != nil || node.Value != "second" {
			t.Errorf("expected head node 'second', got %v (err: %v)", node, err)
		}
		if l.Len() != 0 || !l.IsEmpty() {
			t.Errorf("expected empty list")
		}

		// Remove head from empty list
		_, err = l.RemoveHead()
		if err == nil || err.Error() != "LinkListError: empty list" {
			t.Errorf("expected empty list error when removing head from empty list, got %v", err)
		}
	})

	t.Run("LinkListError custom string output", func(t *testing.T) {
		eUnknown := LinkListError{errCode: 999}
		if eUnknown.Error() != "LinkListError: unknown error" {
			t.Errorf("expected unknown error message, got %s", eUnknown.Error())
		}
	})
}
