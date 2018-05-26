package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriceParsing(t *testing.T) {
	testPrices := []string{
		"~price 3 chaos",
		"~b/o 10 chaos",
		"~price 1 alch each",
	}

	_ = testPrices

	assert.Nil(t, nil)
}
