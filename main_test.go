package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type Input struct {
	Param1 int
	Param2 int
}

type Output = int

type Testcase struct {
	Input  Input
	Output Output
}

func Test(t *testing.T) {
	tests := []Testcase{
		{Input{1, 2}, 3},
		{Input{2, 2}, 4},
		{Input{3, 2}, 5},
	}

	for index, item := range tests {
		t.Run(fmt.Sprintf("test-%d", index+1), func(t *testing.T) {
			result := sum(item.Input.Param1, item.Input.Param2)
			require.Equal(t, item.Output, result, "Test %d failed", index+1)
		})
	}
}
