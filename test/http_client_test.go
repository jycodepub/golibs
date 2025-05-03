package test

import (
	"fmt"
	"github.com/jycodepub/golibs/net"
	"testing"
)

func TestGet(t *testing.T) {
	client := net.NewHttpClient()
	resp, err := client.Get("http://www.google.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Body)
	if !resp.IsOK() {
		t.Fail()
	}
}
