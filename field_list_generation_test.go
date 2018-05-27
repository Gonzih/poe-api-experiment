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
		"+27 to Intelligence": {v: 27, n: "+ to Intelligence"},
		// "101% increased Critical Strike Chance for Spells": {v: 101, n: "% increased Critical Strike Chance for Spells"},
		// "37% increased Mana Regeneration Rate":             {v: 37, n: "% increased Mana Regeneration Rate"},
		// "Adds 6 to 12 Cold Damage to Spells":               {v: 9, n: "Adds 6 to 12 Cold Damage to Spells"},
	}

	for in, out := range data {
		value, name := parseModString(in)
		assert.Equal(t, out.v, value)
		assert.Equal(t, out.n, name)
	}
}
