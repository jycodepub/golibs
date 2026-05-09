package test

import (
	"fmt"
	"testing"

	"github.com/jycodepub/golibs/net"
)

func TestGet(t *testing.T) {
	client := net.NewHttpClient()
	resp, err := client.Get("http://www.google.com", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Body)
	if !resp.IsOK() {
		t.Fail()
	}
}
