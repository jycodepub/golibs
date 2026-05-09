# Template Engine
## Package
- `github.com/jycodepub/golibs/template`

## Type
### 1. `Engine`
Methods:
- `func NewEngine()
- `func (t *Engine) Execute(ctx *Context) error
- `func (t *Engine) ExecuteToFile(ctx *Context, path string) error`

### 2. `Contex`
Fields:
- `template string`
- `data map[string]interface{}`
Methods:
- `func NewContext(templatePath string, dataPath string) (*Context, error)`

## Examples
```go
ctx, err := template.NewContext("template.txt", "data.json")
if err != nil {
	panic(err)
}
engine := template.NewEngine()
err := engine.Execute(ctx)
err := engine.ExecuteToFile(ctx, "output.txt")
```