package test

import (
	"testing"

	"github.com/jycodepub/golibs/template"
)

func TestEngine(t *testing.T) {
	engine := template.NewEngine()
	ctx := &template.Context{}
	data := make(map[string]interface{})
	data["Name"] = "World"
	ctx.Init("Hello {{.Name}}!", data)

	// Execute panics on error, so we can defer recover to catch it
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute panicked: %v", r)
		}
	}()
	err := engine.Execute(ctx)
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}
}
