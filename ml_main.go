package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"

	. "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

func generateTensors(input [][]float64, vecSize int) (ts []tensor.Tensor) {
	cols := make([][]float64, len(input[0]))

	for _, row := range input {
		for i, v := range row {
			cols[i] = append(cols[i], v)
		}
	}

	for _, col := range cols {
		ts = append(ts, tensor.New(tensor.WithBacking(col), tensor.WithShape(vecSize)))
	}

	return
}

//The formula for a straight line is

// y = i0*a0 + i1*a1 + i2*a2 + c

// We want to find an `m` and a `c` that fits the equation well. We'll do it in both float64 and float64 to showcase the extensibility of Gorgonia

func linearRegression(mlInput *MLInput) (func([]float64) (float64, error), error) {
	nInputs := len(mlInput.Fields[0]) - 1
	vecSize := len(mlInput.Fields)
	tensors := generateTensors(mlInput.Fields, vecSize)

	g := NewGraph()

	y := NewVector(g, Float64, WithShape(vecSize), WithName("y"), WithValue(tensors[0]))

	log.Printf("Input len: %d, input[0] len: %d, tensors len: %d", vecSize, nInputs, len(tensors))

	tensors = tensors[1:]

	iss := make([]*Node, nInputs)
	for i := range iss {
		iss[i] = NewVector(g, Float64, WithShape(vecSize), WithName(fmt.Sprintf("i%d", i)), WithValue(tensors[i]))
	}

	sss := make([]*Node, nInputs)
	for i := range sss {
		sss[i] = NewScalar(g, Float64, WithName(fmt.Sprintf("a%d", i)), WithValue(rand.Float64()))
	}

	c := NewScalar(g, Float64, WithName("c"), WithValue(rand.Float64()))

	var expressions []*Node

	for i := range sss {
		expressions = append(expressions, Must(Mul(iss[i], sss[i])))
	}

	expressions = append(expressions, c)

	for {
		if len(expressions) == 1 {
			break
		}

		n1 := expressions[0]
		n2 := expressions[1]

		expressions = append(expressions[2:], Must(Add(n1, n2)))
		// expressions[0] = Must(Add(n1, n2))
	}

	pred := expressions[0]

	se := Must(Square(Must(Sub(pred, y))))
	cost := Must(Mean(se))

	allScalars := append(sss, c)

	_, err := Grad(cost, allScalars...)

	// machine := NewLispMachine(g)  // you can use a LispMachine, but it'll be VERY slow.
	machine := NewTapeMachine(g, BindDualValues(allScalars...))

	defer runtime.GC()
	model := allScalars
	solver := NewVanillaSolver(WithLearnRate(0.001), WithClip(5)) // good idea to clip

	if CUDA {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}
	for i := 0; i < 1000; i++ {
		if err = machine.RunAll(); err != nil {
			fmt.Printf("Error during iteration: %v: %v\n", i, err)
			break
		}

		if err = solver.Step(model); err != nil {
			log.Fatal(err)
		}

		machine.Reset() // Reset is necessary in a loop like this
	}

	var resValues []Value

	for _, s := range allScalars {
		resValues = append(resValues, s.Value())
	}

	// var output strings.Builder
	// output.WriteString("y = ")

	// for i := range sss {
	// 	output.WriteString(fmt.Sprintf("i%d*%3.3f + ", i, sss[i].Value()))
	// }

	// output.WriteString(fmt.Sprintf("%3.3f", c.Value()))

	// log.Println(output.String())

	evaluationFn := func(in []float64) (out float64, err error) {
		in = in[1:]

		if len(in) != nInputs {
			err = fmt.Errorf("Number of inputs is wrong %d != %d", len(in), nInputs)
			return
		}

		for i := range sss {
			out += in[i] * sss[i].Value().Data().(float64)
		}

		out += c.Value().Data().(float64)

		return
	}

	return evaluationFn, nil
}
