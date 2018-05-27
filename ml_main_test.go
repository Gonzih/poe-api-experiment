package main

import (
	"log"
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
		_, err := linearRegression(&MLInput{Fields: input})
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

	evalFn, err := linearRegression(&MLInput{Fields: input})
	assert.Nil(t, err)

	_, err = evalFn([]float32{0, 4, 2})
	assert.NotNil(t, err)

	result, err := evalFn([]float32{0, 4, 2, 1})
	assert.Nil(t, err)
	assert.InEpsilon(t, float32(11), result, 0.1)
}

func TestBasicDoubles(t *testing.T) {
	var input [][]float32

	for i := float32(0); i < 100000; i++ {
		input = append(input, []float32{i*2 + 1, i})
	}

	f, err := linearRegression(&MLInput{Fields: input})
	assert.Nil(t, err)

	result, err := f([]float32{0, 6000})
	assert.Nil(t, err)

	log.Printf("Result %3.3f", result)
	assert.InEpsilon(t, float32(12001), result, 0.01)
}
