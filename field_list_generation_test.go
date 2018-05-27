package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractMod(t *testing.T) {
	data := map[string]struct {
		n string
		v float32
	}{
		"+27 to Intelligence":                              {v: 27, n: "+\\d+ to Intelligence"},
		"101% increased Critical Strike Chance for Spells": {v: 101, n: "\\d+% increased Critical Strike Chance for Spells"},
		"37% increased Mana Regeneration Rate":             {v: 37, n: "\\d+% increased Mana Regeneration Rate"},
		"Adds 6 to 12 Cold Damage to Spells":               {v: 9, n: "Adds \\d+ to \\d+ Cold Damage to Spells"},
		"Arrows Pierce all Targets":                        {v: 1, n: "Arrows Pierce all Targets"},
	}

	for in, out := range data {
		value, name := parseModString(in)
		assert.Equal(t, out.v, value)
		assert.Equal(t, out.n, name)
	}
}
