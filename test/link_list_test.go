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


 	fmt.Println(list.Get(0).Value)
 	fmt.Println(list.Get(1).Value)
 	fmt.Println(list.Get(2).Value)

	for i := range list.Len() {
		fmt.Println(list.Get(i).Value)
	}

	for !list.IsEmpty() {
		fmt.Println(list.PopHead().Value)
	}
}