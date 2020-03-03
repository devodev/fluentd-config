package config

import (
	"fmt"
	"strings"
)

// Element .
type Element interface {
	Print() string
}

// Document .
type Document struct {
	Elements []Element
}

// Print .
func (d *Document) Print() string {
	b := strings.Builder{}

	for _, e := range d.Elements {
		b.WriteString(e.Print())
	}

	return b.String()
}

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

	// opening bracket
	if p.Pattern == "" {
		b.WriteString(fmt.Sprintf("<%s>\n", p.Name))
	} else {
		b.WriteString(fmt.Sprintf("<%s %s>\n", p.Name, p.Pattern))
	}

	// parameters
	for _, param := range p.Parameters {
		b.WriteString(fmt.Sprintf("  %s %s\n", param.Key, param.Value))
	}

	// blocks
	for _, block := range p.Blocks {
		// block opening bracket
		if block.Pattern == "" {
			b.WriteString(fmt.Sprintf("  <%s>\n", block.Name))
		} else {
			b.WriteString(fmt.Sprintf("  <%s %s>\n", block.Name, block.Pattern))
		}
		// block parameters
		for _, param := range block.Parameters {
			b.WriteString(fmt.Sprintf("    %s %s\n", param.Key, param.Value))
		}
		// block closing bracket
		b.WriteString(fmt.Sprintf("  </%s>\n", block.Name))
	}

	// closing bracket
	b.WriteString(fmt.Sprintf("</%s>\n", p.Name))

	return b.String()
}

// Block .
type Block struct {
	Name       string
	Pattern    string
	Parameters []Parameter
}

// Parameter .
type Parameter struct {
	Key   string
	Value string
}

// Include .
type Include struct {
	Value string
}

// Print .
func (i *Include) Print() string {
	return fmt.Sprintf("@include %s\n", i.Value)
}
