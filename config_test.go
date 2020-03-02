package config

import (
	"fmt"
	"testing"
)

func TestConfigPrint(t *testing.T) {

	cases := []struct {
		Name   string
		Plugin *Plugin
		Want   string
	}{
		{
			Name: "empty source",
			Plugin: &Plugin{
				Name:       "source",
				Pattern:    "",
				Parameters: []Parameter{},
				Blocks:     []Block{},
			},
			Want: "<source>\n</source>\n",
		},
		{
			Name: "source with parameters",
			Plugin: &Plugin{
				Name:    "source",
				Pattern: "",
				Parameters: []Parameter{
					Parameter{
						Key:   "@type",
						Value: "forward",
					},
					Parameter{
						Key:   "port",
						Value: "24224",
					},
				},
				Blocks: []Block{},
			},
			Want: "<source>\n" +
				"  @type forward\n" +
				"  port 24224\n" +
				"</source>\n",
		},
	}

	for idx, c := range cases {
		t.Run(fmt.Sprintf("%d. %s", idx, c.Name), func(t *testing.T) {
			got := c.Plugin.Print()
			if got != c.Want {
				t.Errorf("\ngot:\n%v\nbut want:\n%v\n", got, c.Want)
			}
		})
	}
}
