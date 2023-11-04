package main

import (
	"fmt"
	"testing"

	"github.com/MaoScut/best-model-assemble-seq/structs"
)

func TestGetAll(t *testing.T) {
	input := []structs.Unit{
		{},
		{},
	}
	input2 := input[2:]
	fmt.Println(input2)
	result := remove(input, 1)
	fmt.Println(len(result))
}

func TestGetAllSolution(t *testing.T) {
	input := []structs.Unit{
		{
			Name: "arms1",
		},
		{
			Name: "arms2",
		},
		{
			Name: "arms3",
		},
		{
			Name: "arms4",
		},
		{
			Name: "arms5",
		},
		{
			Name: "arms6",
		},
		{
			Name: "arms7",
		},
	}
	all := getAllSolutions(input)
	// expectedCount := 5040
	// if len(all) != expectedCount {
	// 	t.Fatal("should have expected count solution", len(all), expectedCount)
	// }
	for _, s := range all {
		t.Log(s.String())
	}
}
