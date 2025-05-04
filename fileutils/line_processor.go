package fileutils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type LineProcessor struct {
	file *os.File
}

func (p *LineProcessor) Open(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	p.file = f
	return nil
}

func (p *LineProcessor) Close() {
	_ = p.file.Close()
}

type Accumulator interface {
	Accumulate(o string)
}

func (p *LineProcessor) Process(f func(string) (string, error), a Accumulator) (int, error) {
	cnt := 0
	r := bufio.NewReader(p.file)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		} else {
			o, err := f(line)
			if err != nil {
				continue
			}
			if a != nil {
				a.Accumulate(o)
			}
			cnt += 1
		}
	}
	return cnt, nil
}
