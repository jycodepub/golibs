package dstruct

import (
	"testing"
)

func TestStack(t *testing.T) {
	t.Run("NewStack and IsEmpty", func(t *testing.T) {
		s := NewStack[int](3)
		if !s.IsEmpty() {
			t.Errorf("expected new stack to be empty")
		}
	})

	t.Run("Push and Pop basic behavior", func(t *testing.T) {
		s := NewStack[string](2)

		err := s.Push("first")
		if err != nil {
			t.Fatalf("unexpected push error: %v", err)
		}
		if s.IsEmpty() {
			t.Errorf("expected stack to not be empty after push")
		}

		err = s.Push("second")
		if err != nil {
			t.Fatalf("unexpected push error: %v", err)
		}

		// Pop second item (LIFO)
		val, popErr := s.Pop()
		if popErr != nil {
			t.Fatalf("unexpected pop error: %v", popErr)
		}
		if val != "second" {
			t.Errorf("expected 'second', got '%s'", val)
		}

		// Pop first item
		val, popErr = s.Pop()
		if popErr != nil {
			t.Fatalf("unexpected pop error: %v", popErr)
		}
		if val != "first" {
			t.Errorf("expected 'first', got '%s'", val)
		}

		if !s.IsEmpty() {
			t.Errorf("expected stack to be empty after popping all elements")
		}
	})

	t.Run("StackFull error", func(t *testing.T) {
		s := NewStack[int](2)

		_ = s.Push(10)
		_ = s.Push(20)

		err := s.Push(30)
		if err == nil {
			t.Fatalf("expected StackFull error, got nil")
		}
		if err.Error() != "stack is full" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("StackEmpty error", func(t *testing.T) {
		s := NewStack[int](2)

		val, err := s.Pop()
		if err == nil {
			t.Fatalf("expected StackEmpty error, got nil")
		}
		if val != 0 {
			t.Errorf("expected zero value 0, got %d", val)
		}
		if err.Error() != "stack is empty" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})
}
