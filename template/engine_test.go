package template

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNewContext_Success(t *testing.T) {
	tmplPath := filepath.Join("..", "test", "sample.tmpl")
	dataPath := filepath.Join("..", "test", "sample_data.json")

	ctx, err := NewContext(tmplPath, dataPath)
	if err != nil {
		t.Fatalf("unexpected error from NewContext: %v", err)
	}

	if ctx.template != "Hello {{.Name}}, welcome to {{.Platform}}!\n" {
		t.Errorf("unexpected template loaded: %q", ctx.template)
	}

	if ctx.data["Name"] != "Alice" || ctx.data["Platform"] != "GoLibs" {
		t.Errorf("unexpected data loaded: %+v", ctx.data)
	}
}

func TestNewContext_FileErrors(t *testing.T) {
	tmplPath := filepath.Join("..", "test", "sample.tmpl")
	dataPath := filepath.Join("..", "test", "sample_data.json")

	// Missing template file
	_, err := NewContext("non_existent.tmpl", dataPath)
	if err == nil {
		t.Errorf("expected error for non-existent template file, got nil")
	}

	// Missing data file
	_, err = NewContext(tmplPath, "non_existent.json")
	if err == nil {
		t.Errorf("expected error for non-existent data file, got nil")
	}

	// Invalid JSON data file
	_, err = NewContext(tmplPath, tmplPath)
	if err == nil {
		t.Errorf("expected error for invalid JSON data file, got nil")
	}
}

func TestContext_Init(t *testing.T) {
	ctx := &Context{}
	tmplStr := "User: {{.User}}"
	dataMap := map[string]interface{}{"User": "Bob"}

	ctx.Init(tmplStr, dataMap)

	if ctx.template != tmplStr {
		t.Errorf("expected template %q, got %q", tmplStr, ctx.template)
	}
	if ctx.data["User"] != "Bob" {
		t.Errorf("expected User Bob, got %v", ctx.data["User"])
	}
}

func TestEngine_Execute(t *testing.T) {
	ctx := &Context{}
	ctx.Init("Hello {{.Name}}!", map[string]interface{}{"Name": "Charlie"})

	// Intercept stdout to test Execute(ctx)
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	engine := NewEngine()
	execErr := engine.Execute(ctx)

	w.Close()
	os.Stdout = oldStdout

	if execErr != nil {
		t.Fatalf("unexpected error executing engine: %v", execErr)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	expected := "Hello Charlie!"
	if buf.String() != expected {
		t.Errorf("expected output %q, got %q", expected, buf.String())
	}
}

func TestEngine_ExecuteToFile(t *testing.T) {
	tmplPath := filepath.Join("..", "test", "sample.tmpl")
	dataPath := filepath.Join("..", "test", "sample_data.json")

	ctx, err := NewContext(tmplPath, dataPath)
	if err != nil {
		t.Fatalf("failed to create context: %v", err)
	}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "output.txt")

	engine := NewEngine()
	err = engine.ExecuteToFile(ctx, outputPath)
	if err != nil {
		t.Fatalf("unexpected error executing to file: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	expected := "Hello Alice, welcome to GoLibs!\n"
	if string(content) != expected {
		t.Errorf("expected content %q, got %q", expected, string(content))
	}
}

func TestEngine_ParseAndExecuteErrors(t *testing.T) {
	// Invalid template syntax
	ctxInvalid := &Context{}
	ctxInvalid.Init("Hello {{.Name", map[string]interface{}{"Name": "Test"})

	err := execute(ctxInvalid, &bytes.Buffer{})
	if err == nil {
		t.Errorf("expected template parse error, got nil")
	}
}
