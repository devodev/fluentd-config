package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestPluginIncludePrint(t *testing.T) {

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

func TestDocumentPrint(t *testing.T) {

	stubInclude := &Include{Value: "file.conf"}
	stubSourceEmpty := &Plugin{Name: "source"}
	stubSourceWithParams := &Plugin{
		Name: "source",
		Parameters: []Parameter{
			Parameter{Key: "@type", Value: "forward"},
			Parameter{Key: "port", Value: "24224"},
		},
	}

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
		{
			Name:     "document with one empty source",
			Document: &Document{[]Element{stubSourceEmpty}},
			Want:     "<source>\n</source>\n",
		},
		{
			Name:     "document with one source with parameters",
			Document: &Document{[]Element{stubSourceWithParams}},
			Want: "<source>\n" +
				"  @type forward\n" +
				"  port 24224\n" +
				"</source>\n",
		},
		{
			Name:     "document with one source with parameters and one include",
			Document: &Document{[]Element{stubSourceWithParams, stubInclude}},
			Want: "<source>\n" +
				"  @type forward\n" +
				"  port 24224\n" +
				"</source>\n" +
				"@include file.conf\n",
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

func TestDocumentJSONDecode(t *testing.T) {

	cases := []struct {
		Name string
		JSON string
		Want *Document
	}{
		{Name: "emtpy document", JSON: `{"elements": []}`, Want: &Document{Elements: []Element{}}},
		{
			Name: "document with one include",
			JSON: `{"elements": [{"type": "include", "data": {"value": "file.conf"}}]}`,
			Want: &Document{
				Elements: []Element{
					&Include{Value: "file.conf"},
				},
			},
		},
		{
			Name: "document with one empty source",
			JSON: `{"elements":[{"type":"plugin","data":{"name":"source"}}]}`,
			Want: &Document{
				Elements: []Element{
					&Plugin{Name: "source"},
				},
			},
		},
		{
			Name: "document with one source with parameters",
			JSON: `{"elements":[{"type":"plugin","data":{"name":"source","parameters":[{"key":"@type","value":"forward"}]}}]}`,
			Want: &Document{
				Elements: []Element{
					&Plugin{Name: "source", Parameters: []Parameter{Parameter{Key: "@type", Value: "forward"}}},
				},
			},
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d. %s", idx, c.Name), func(t *testing.T) {
			var got Document
			b := bytes.NewBuffer([]byte(c.JSON))
			err := json.NewDecoder(b).Decode(&got)
			if err != nil {
				t.Fatalf("error occured while decoding JSON: %v", err)
			}

			if !reflect.DeepEqual(&got, c.Want) {
				t.Errorf("\ngot:\n%v\nbut want:\n%v\n", &got, c.Want)
			}
		})
	}
}
