package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapFrameType(t *testing.T) {
	assert.Equal(t, "normal", mapFrameType(0))
	assert.Equal(t, "magic", mapFrameType(1))
	assert.Equal(t, "rare", mapFrameType(2))
	assert.Equal(t, "unique", mapFrameType(3))
	assert.Equal(t, "gem", mapFrameType(4))
	assert.Equal(t, "currency", mapFrameType(5))
	assert.Equal(t, "divination card", mapFrameType(6))
	assert.Equal(t, "quest item", mapFrameType(7))
	assert.Equal(t, "prophecy", mapFrameType(8))
	assert.Equal(t, "relic", mapFrameType(9))
	assert.Equal(t, "", mapFrameType(999))
}

func TestPriceParsing(t *testing.T) {
	testPrices := map[string]float64{
		"~price 3 chaos":      42,
		"~b/o 10 chaos":       140,
		"~price 1 chaos each": 14,
		"~price 1 exa":        700.000,
		"~b/o 1 chisel":       3.500,
		"~b/o 1 chaos":        14.000,
		"~b/o 1 fuse":         7.000,
		"~b/o 1 alch":         4.000,
		"~price 10 exa":       7000.000,
		"~b/o 2 chrom":        2.000,
		"~b/o 1 chrom":        1.000,
	}

	for priceS, parsed := range testPrices {
		price, ok := parsePriceInChrom(priceS)
		assert.True(t, ok)
		assert.Equal(t, parsed, price)
	}
}

func TestMLInputPersistancy(t *testing.T) {
	fname := "/tmp/test-ml-input.bin"

	in := &MLInput{Fields: [][]float64{
		[]float64{1, 2, 3},
		[]float64{4, 5, 6},
		[]float64{7, 8, 9},
	}, FName: fname}

	err := in.Save()
	assert.Nil(t, err)

	in2 := &MLInput{FName: fname}

	err = in2.Load()
	assert.Nil(t, err)

	assert.Equal(t, len(in.Fields), len(in2.Fields))
	assert.Equal(t, len(in.Fields[0]), len(in2.Fields[0]))

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			assert.Equal(t, in.Fields[i][j], in2.Fields[i][j])
		}
	}
}
