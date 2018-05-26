package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"

	. "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

func generateTensors(input [][]float32, vecSize int) (ts []tensor.Tensor) {
	cols := make([][]float32, len(input[0]))

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

// We want to find an `m` and a `c` that fits the equation well. We'll do it in both float32 and float32 to showcase the extensibility of Gorgonia

func linearRegression(input [][]float32) error {
	nInputs := len(input[0]) - 1
	vecSize := len(input)
	tensors := generateTensors(input, vecSize)

	g := NewGraph()

	y := NewVector(g, Float32, WithShape(vecSize), WithName("y"), WithValue(tensors[0]))

	log.Printf("Input len: %d, input[0] len: %d, tensors len: %d", len(input), len(input[0]), len(tensors))

	tensors = tensors[1:]

	iss := make([]*Node, nInputs)
	for i := range iss {
		iss[i] = NewVector(g, Float32, WithShape(vecSize), WithName(fmt.Sprintf("i%d", i)), WithValue(tensors[i]))
	}

	sss := make([]*Node, nInputs)
	for i := range sss {
		sss[i] = NewScalar(g, Float32, WithName(fmt.Sprintf("a%d", i)), WithValue(rand.Float32()))
	}

	c := NewScalar(g, Float32, WithName("c"), WithValue(rand.Float32()))

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
	for i := 0; i < 500; i++ {
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

	log.Printf("y = i0*%3.3f + i1*%3.3f + i2*%3.3f + %3.3f\n",
		resValues[0],
		resValues[1],
		resValues[2],
		resValues[3],
	)

	return nil
}

func MLMain(input [][]float32) {
	linearRegression(input)
}
