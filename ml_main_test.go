package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicExecution(t *testing.T) {
	inputs := [][][]float64{
		[][]float64{
			[]float64{1, 2, 3, 4},
			[]float64{1, 2, 3, 5},
			[]float64{1, 1, 3, 5},
			[]float64{1, 2, 1, 5},
			[]float64{0, 2, 1, 5},
		},
		[][]float64{
			[]float64{1, 2, 3, 4, 1},
			[]float64{1, 2, 3, 5, 1},
			[]float64{1, 1, 3, 5, 1},
			[]float64{1, 2, 1, 5, 1},
			[]float64{0, 2, 1, 5, 1},
		},
		[][]float64{
			[]float64{1, 2, 3, 4, 1, 2, 3, 5},
			[]float64{1, 2, 3, 5, 1, 2, 3, 5},
			[]float64{1, 1, 3, 5, 1, 2, 3, 5},
			[]float64{1, 2, 1, 5, 1, 2, 3, 5},
			[]float64{0, 2, 1, 5, 1, 2, 3, 5},
		},
	}

	for _, input := range inputs {
		_, err := linearRegression(&MLInput{Fields: input})
		assert.Nil(t, err)
	}
}

func TestBasicEvaluation(t *testing.T) {
	input := [][]float64{
		[]float64{16, 5, 3, 4},
		[]float64{12, 3, 3, 5},
		[]float64{15, 4, 3, 5},
		[]float64{10, 2, 1, 5},
		[]float64{20, 6, 1, 5},
	}

	evalFn, err := linearRegression(&MLInput{Fields: input})
	assert.Nil(t, err)

	_, err = evalFn([]float64{0, 4, 2})
	assert.NotNil(t, err)

	result, err := evalFn([]float64{0, 4, 2, 1})
	assert.Nil(t, err)
	assert.InEpsilon(t, float64(11), result, 0.1)
}

func TestBasicDoubles(t *testing.T) {
	var input [][]float64

	for i := float64(0); i < 100000; i++ {
		input = append(input, []float64{i*2 + 1, i})
	}

	f, err := linearRegression(&MLInput{Fields: input})
	assert.Nil(t, err)

	result, err := f([]float64{0, 6000})
	assert.Nil(t, err)

	log.Printf("Result %3.3f", result)
	assert.InEpsilon(t, float64(12001), result, 0.01)
}
