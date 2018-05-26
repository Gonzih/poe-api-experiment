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
		_, err := linearRegression(input)
		assert.Nil(t, err)
	}
}

func TestBasicEvaluation(t *testing.T) {
	input := [][]float32{
		[]float32{16, 5, 3, 4},
		[]float32{12, 3, 3, 5},
		[]float32{15, 4, 3, 5},
		[]float32{10, 2, 1, 5},
		[]float32{20, 6, 1, 5},
	}

	evalFn, err := linearRegression(input)
	assert.Nil(t, err)

	_, err = evalFn([]float32{0, 4, 2})
	assert.NotNil(t, err)

	result, err := evalFn([]float32{0, 4, 2, 1})
	assert.Nil(t, err)
	assert.Equal(t, float32(11.142875), result)
}
