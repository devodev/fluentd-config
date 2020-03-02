package config

import "testing"

func TestConfig(t *testing.T) {

	b := &Plugin{
		Name:       "source",
		Pattern:    "",
		Parameters: []Parameter{},
		Blocks:     []Block{},
	}

	got := b.Print()
	want := "<source>\n</source>\n"

	if got != want {
		t.Errorf("\ngot:\n%q\nbut want:\n%q\n", got, want)
	}
}
