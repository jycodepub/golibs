package test

import (
	"fmt"
	"github.com/jycodepub/golibs/fileutils"
	"testing"
)

func TestProcess(t *testing.T) {
	var p fileutils.LineProcessor
	err := p.Open("line_processor_test.go")
	if err != nil {
		panic(err)
	}
	defer p.Close()
	cnt, err := p.Process(processor, &myAccumulator{})
	if err != nil {
		panic(err)
	}
	fmt.Println(cnt)
}

func processor(line string) string {
	fmt.Print("Processed: ", line)
	return line
}

type myAccumulator struct {
}

func (a *myAccumulator) Accumulate(o string) {
	fmt.Print("Accumulated: ", o)
}
