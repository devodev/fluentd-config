package config

import (
	"fmt"
	"testing"
)

func TestConfigPrint(t *testing.T) {

	cases := []struct {
		Name    string
		Plugin  *Plugin
		Include *Include
		Want    string
	}{
		{
			Name:   "empty source",
			Plugin: &Plugin{Name: "source"},
			Want:   "<source>\n</source>\n",
		},
		{
			Name: "source with parameters",
			Plugin: &Plugin{
				Name: "source",
				Parameters: []Parameter{
					Parameter{Key: "@type", Value: "forward"},
					Parameter{Key: "port", Value: "24224"},
				},
			},
			Want: "<source>\n" +
				"  @type forward\n" +
				"  port 24224\n" +
				"</source>\n",
		},
		{
			Name:   "filter with pattern",
			Plugin: &Plugin{Name: "filter", Pattern: "**"},
			Want: "<filter **>\n" +
				"</filter>\n",
		},
		{
			Name: "filter with parameters and a block",
			Plugin: &Plugin{
				Name:    "filter",
				Pattern: "myapp.access",
				Parameters: []Parameter{
					Parameter{Key: "@type", Value: "record_transformer"},
				},
				Blocks: []Block{
					Block{
						Name:    "record",
						Pattern: "",
						Parameters: []Parameter{
							Parameter{Key: "host_param", Value: "\"#{Socket.gethostname}\""},
						},
					},
				},
			},
			Want: "<filter myapp.access>\n" +
				"  @type record_transformer\n" +
				"  <record>\n" +
				"    host_param \"#{Socket.gethostname}\"\n" +
				"  </record>\n" +
				"</filter>\n",
		},
		{
			Name: "include",
			Include: &Include{
				Value: "file.conf",
			},
			Want: "@include file.conf\n",
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d. %s", idx, c.Name), func(t *testing.T) {
			var got string
			if c.Include == nil {
				got = c.Plugin.Print()
			} else {
				got = c.Include.Print()
			}

			if got != c.Want {
				t.Errorf("\ngot:\n%v\nbut want:\n%v\n", got, c.Want)
			}
		})
	}
}

func TestDocument(t *testing.T) {

	stubInclude := &Include{Value: "file.conf"}

	cases := []struct {
		Name     string
		Document *Document
		Want     string
	}{
		{
			Name: "emtpy document", Document: &Document{}, Want: "",
		},
		{
			Name:     "document with one include",
			Document: &Document{[]Element{stubInclude}},
			Want:     "@include file.conf\n",
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d. %s", idx, c.Name), func(t *testing.T) {
			got := c.Document.Print()

			if got != c.Want {
				t.Errorf("\ngot:\n%v\nbut want:\n%v\n", got, c.Want)
			}
		})
	}
}
