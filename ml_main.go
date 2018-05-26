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
	cols := make([][]float32, vecSize)

	for _, row := range input {
		for i, v := range row {
			cols[i] = append(cols[i], v)
		}
	}

	// log.Printf("input: %v", input)
	// log.Printf("cols: %v", cols)

	for _, col := range cols {
		ts = append(ts, tensor.New(tensor.WithBacking(col), tensor.WithShape(vecSize)))
	}

	return
}

//The formula for a straight line is

// y = i0*a0 + i1*a1 + i2*a2 + c

// We want to find an `m` and a `c` that fits the equation well. We'll do it in both float32 and float32 to showcase the extensibility of Gorgonia

func linearRegression(input [][]float32) {
	vecSize := len(input)
	tensors := generateTensors(input, vecSize)

	g := NewGraph()

	y := NewVector(g, Float32, WithShape(vecSize), WithName("y"), WithValue(tensors[0]))

	i1 := NewVector(g, Float32, WithShape(vecSize), WithName("i1"), WithValue(tensors[1]))
	i2 := NewVector(g, Float32, WithShape(vecSize), WithName("i2"), WithValue(tensors[2]))
	i3 := NewVector(g, Float32, WithShape(vecSize), WithName("i3"), WithValue(tensors[3]))

	a1 := NewScalar(g, Float32, WithName("a1"), WithValue(rand.Float32()))
	a2 := NewScalar(g, Float32, WithName("a2"), WithValue(rand.Float32()))
	a3 := NewScalar(g, Float32, WithName("a3"), WithValue(rand.Float32()))

	c := NewScalar(g, Float32, WithName("c"), WithValue(rand.Float32()))

	pred := Must(Add(
		Must(Add(
			Must(Mul(i1, a1)),
			Must(Mul(i2, a2)),
		)),
		Must(Add(
			Must(Mul(i3, a3)),
			c,
		)),
	))

	se := Must(Square(Must(Sub(pred, y))))
	cost := Must(Mean(se))

	_, err := Grad(cost, a1, a2, a3, c)

	// machine := NewLispMachine(g)  // you can use a LispMachine, but it'll be VERY slow.
	machine := NewTapeMachine(g, BindDualValues(a1, a2, a3, c))

	defer runtime.GC()
	model := Nodes{a1, a2, a3, c}
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

	log.Printf("y = i0*%3.3f + i1*%3.3f + i2*%3.3f + %3.3f\n", a1.Value(), a2.Value(), a3.Value(), c.Value())
}

func MLMain(input [][]float32) {
	linearRegression(input)
}
