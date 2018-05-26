package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicExecution(t *testing.T) {
	inputs := [][][]float32{
		[][]float32{
			[]float32{1, 2, 3, 4},
			[]float32{1, 2, 3, 5},
			[]float32{1, 1, 3, 5},
			[]float32{1, 2, 1, 5},
			[]float32{0, 2, 1, 5},
		},
		[][]float32{
			[]float32{1, 2, 3, 4, 1},
			[]float32{1, 2, 3, 5, 1},
			[]float32{1, 1, 3, 5, 1},
			[]float32{1, 2, 1, 5, 1},
			[]float32{0, 2, 1, 5, 1},
		},
		[][]float32{
			[]float32{1, 2, 3, 4, 1, 2, 3, 5},
			[]float32{1, 2, 3, 5, 1, 2, 3, 5},
			[]float32{1, 1, 3, 5, 1, 2, 3, 5},
			[]float32{1, 2, 1, 5, 1, 2, 3, 5},
			[]float32{0, 2, 1, 5, 1, 2, 3, 5},
		},
	}

	for _, input := range inputs {
		err := linearRegression(input)
		assert.Nil(t, err)
	}
}
