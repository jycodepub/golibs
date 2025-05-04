package test

import (
	"fmt"
	"github.com/jycodepub/golibs/dstruct"
	"testing"
)

func TestStack(t *testing.T) {
	s := dstruct.NewStack[int](2)
	err := s.Push(1)
	if err != nil {
		t.Error("Push 1 occurred error, expect no error")
	}
	err = s.Push(2)
	if err != nil {
		t.Error("Push 2 occurred error, expect no error")
	}
	err = s.Push(3)
	if err.Error() != "stack is full" {
		t.Error("Push 3 no error, expect StackFull error")
	}
	e, _ := s.Pop()
	if e != 2 {
		t.Errorf("Pop %d, expect 2", e)
	}
	e, _ = s.Pop()
	if e != 1 {
		t.Errorf("Pop %d, expect 1", e)
	}
	e, err2 := s.Pop()
	if err2.Error() != "stack is empty" {
		t.Error("Pop expect stack is empty")
	}
}

func ExampleStack() {
	s := dstruct.NewStack[int](10)
	for i := 0; i < 10; i++ {
		_ = s.Push(i)
	}
	for !s.IsEmpty() {
		e, _ := s.Pop()
		fmt.Println(e)
	}
}
