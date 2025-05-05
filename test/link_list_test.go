package test

import (
	"fmt"
	"testing"

	"github.com/jycodepub/golibs/dstruct"
)

func TestList(t *testing.T) {
	list := dstruct.LinkList[int]{}
	list.Append(&dstruct.Node[int]{Value: 1})
	list.Append(&dstruct.Node[int]{Value: 2})
	list.Append(&dstruct.Node[int]{Value: 3})

	i := list.Head
	for i != nil {
		fmt.Println(i.Value)
		i = i.Next
	}

	for !list.IsEmpty() {
		fmt.Println(list.PopHead().Value)
	}
}