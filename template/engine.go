// Package template provides a template engine
package template

import (
	"io"
	"os"
	"text/template"
)

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (t *Engine) Execute(ctx *Context) error {
	return execute(ctx, os.Stdout)
}

func (t *Engine) ExecuteToFile(ctx *Context, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return execute(ctx, file)
}

func execute(ctx *Context, writer io.Writer) error {
	tmpl, err := template.New("template").Parse(ctx.template)
	if err != nil {
		return err
	}
	err = tmpl.Execute(writer, ctx.data)
	if err != nil {
		return err
	}
	return nil
}
