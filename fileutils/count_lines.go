package fileutils

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var sources = []string{".js", ".ts", ".java", ".kt", ".kts", ".go"}

func CountLines(dir string) int {
	adir := dir
	if !filepath.IsAbs(dir) {
		adir, _ = filepath.Abs(dir)
	}

	count := 0
	if err := filepath.Walk(adir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && isSource(filepath.Ext(path)) {
			count += fileLOC(path)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return count
}

func isSource(ext string) bool {
	for _, s := range sources {
		if s == ext {
			return true
		}
	}
	return false
}

func fileLOC(path string) int {
	f, err := os.Open(path)
	if err != nil {
		log.Println("Failed to open file")
		return 0
	}
	defer f.Close()

	count := 0
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		line := strings.TrimSpace(string(b))
		if err == nil {
			if line != "" && !isComment(line) {
				count++
			}
		} else if err == io.EOF {
			break
		} else {
			log.Println("Failed to read file")
		}
	}
	fmt.Println("-", path, "-> LOC:", count)
	return count
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "//")
}
