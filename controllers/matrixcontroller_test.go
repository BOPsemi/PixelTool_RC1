package controllers

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func Test_NewMatrixController(t *testing.T) {
	obj := NewMatrixController()

	assert.NotNil(t, obj)
}

func Test_EvalLinearMatrix(t *testing.T) {
	obj := NewMatrixController()

	rgb := []float64{1.0, 1.0, 1.0, 1.0}
	elm := []float64{1.0, 0.0, 0.0, 0.0, 0.0, 0.0}

	result := obj.EvalLinearMatrix(elm, rgb)

	fmt.Println(result)
}
