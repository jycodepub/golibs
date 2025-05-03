package fileutils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Compact(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer destFile.Close()
	reader := bufio.NewReader(srcFile)
	writer := bufio.NewWriter(destFile)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if len(line) != 0 {
			_, err := writer.WriteString(line)
			if err != nil {
				return err
			}
			writer.Flush()
		}
	}
	fmt.Printf("Created file: %s\n", dest)
	return nil
}
