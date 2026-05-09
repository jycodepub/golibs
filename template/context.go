package template

import (
	"encoding/json"
	"fmt"
	"os"
)

type Context struct {
	template string
	data     map[string]interface{}
}

func NewContext(templatePath string, dataPath string) (*Context, error) {
	bytes, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}
	template := string(bytes)

	bytes, err = os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &Context{
		template: template,
		data:     data,
	}, nil
}

func (c *Context) Init(template string, data map[string]interface{}) {
	c.template = template
	c.data = data
}
