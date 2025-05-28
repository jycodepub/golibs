package test

import (
	"github.com/jycodepub/golibs/config"
	"testing"
)

type Values struct {
	Key  string   `json:"key"`
	List []string `json:"list"`
	Num  int      `json:"num"`
}

func TestConfig(t *testing.T) {
	var values Values
	err := config.Load("config.json", &values)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(values)
	if values.Num != 1 {
		t.Errorf("num should be 1, but got %d", values.Num)
	}
	if values.Key != "value" {
		t.Errorf("key should be key, but got %s", values.Key)
	}
	if len(values.List) != 2 {
		t.Errorf("list should be 2, but got %d", len(values.List))
	}
}
