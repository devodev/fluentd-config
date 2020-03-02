package config

import (
	"fmt"
	"strings"
)

// Plugin .
type Plugin struct {
	Name       string
	Pattern    string
	Parameters []Parameter
	Blocks     []Block
}

// Print .
func (p *Plugin) Print() string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("<%s>\n", p.Name))

	for _, param := range p.Parameters {
		b.WriteString(fmt.Sprintf("  %s %s\n", param.Key, param.Value))
	}

	b.WriteString(fmt.Sprintf("</%s>\n", p.Name))

	return b.String()
}

// Block .
type Block struct{}

// Parameter .
type Parameter struct {
	Key   string
	Value string
}
