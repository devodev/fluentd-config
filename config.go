package config

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Element .
type Element interface {
	Print() string
}

// Document .
type Document struct {
	Elements []Element `json:"elements"`
}

// UnmarshalJSON .
func (d *Document) UnmarshalJSON(b []byte) error {
	var doc DocumentJSON
	if err := json.Unmarshal(b, &doc); err != nil {
		return err
	}
	elements := []Element{}
	for _, e := range doc.Elements {
		switch e.Type {
		case "plugin":
			var elem Plugin
			if err := json.Unmarshal(e.Data, &elem); err != nil {
				return err
			}
			elements = append(elements, &elem)
		case "include":
			var elem Include
			if err := json.Unmarshal(e.Data, &elem); err != nil {
				return err
			}
			elements = append(elements, &elem)
		default:
			return fmt.Errorf("unsupported element type: %s", e.Type)
		}
	}
	*d = Document{elements}
	return nil
}

// Print .
func (d *Document) Print() string {
	b := strings.Builder{}

	for _, e := range d.Elements {
		b.WriteString(e.Print())
	}

	return b.String()
}

// DocumentJSON .
type DocumentJSON struct {
	Elements []ElementJSON `json:"elements"`
}

// ElementJSON .
type ElementJSON struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Include .
type Include struct {
	Value string `json:"value"`
}

// Print .
func (i *Include) Print() string {
	return fmt.Sprintf("@include %s\n", i.Value)
}

// Plugin .
type Plugin struct {
	Name       string      `json:"name"`
	Pattern    string      `json:"pattern"`
	Parameters []Parameter `json:"parameters"`
	Blocks     []Block     `json:"blocks"`
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
	Name       string      `json:"name"`
	Pattern    string      `json:"pattern"`
	Parameters []Parameter `json:"parameters"`
}

// Parameter .
type Parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
