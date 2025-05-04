package test

import (
	"encoding/json"
	"fmt"
	"github.com/jycodepub/golibs/fileutils"
	"strings"
	"testing"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestProcess(t *testing.T) {
	var p fileutils.LineProcessor
	err := p.Open("users.csv")
	if err != nil {
		panic(err)
	}
	defer p.Close()
	ac := myAccumulator{}
	_, err = p.Process(processor, &ac)
	if err != nil {
		panic(err)
	}
	fmt.Println(ac.lineCnt)
}

func processor(line string) (string, error) {
	tokens := strings.Split(line, ",")
	u := User{
		Username: tokens[0],
		Password: tokens[1],
	}
	bytes, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type myAccumulator struct {
	lineCnt int
}

func (a *myAccumulator) Accumulate(o string) {
	var u User
	_ = json.Unmarshal([]byte(o), &u)
	fmt.Println("Processed: ", u)
	a.lineCnt += 1
}
