package main

import (
	"fmt"
	"testing"

	"gorgonia.org/tensor"
)

func TestTensor(t *testing.T) {
	d := tensor.New(tensor.WithBacking([]float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}), tensor.WithShape(2, 5))
	fmt.Println(d)
	fmt.Println(d.Sum(0))
	fmt.Println(d.Sum(1))
	fmt.Println(d.Size())
	fmt.Println(d.Dims())
	fmt.Println(d.Strides())

}
