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

	var node *dstruct.Node[int]
	node, _ = list.Get(0)
 	fmt.Println(node.Value)
	node, _ = list.Get(1)
 	fmt.Println(node.Value)
	node, _ = list.Get(2)
 	fmt.Println(node.Value)

	for i := range list.Len() {
		node, _ := list.Get(i)
		fmt.Println(node.Value)
	}

	for !list.IsEmpty() {
		node, _ := list.RemoveHead()
		fmt.Println(node.Value)
	}

	_, err := list.RemoveHead()
	fmt.Println(err.Error())

	_, err = list.Get(3)
	fmt.Println(err.Error())
}