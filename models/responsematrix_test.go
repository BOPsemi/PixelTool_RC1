package models

import "testing"

import "github.com/stretchr/testify/assert"
import "fmt"

func Test_ElementSerializer(t *testing.T) {
	// XYZ Matrix
	obj1 := &XYZMatrix{
		X: 1.0,
		Y: 2.0,
		Z: 3.0,
	}

	result1 := ElementSerializer(obj1)
	assert.EqualValues(t, 3, len(result1))
	fmt.Println("XYZMatrix", result1)

	// RGB Matrix
	obj2 := &RGBMatrix{
		R: 3.0,
		G: 2.0,
		B: 1.0,
	}

	result2 := ElementSerializer(obj2)
	assert.EqualValues(t, 3, len(result2))
	fmt.Println("RGBMatrix", result2)

	// Lab Matrix
	obj3 := &LabMatrix{
		L: 2.0,
		A: 1.0,
		B: 3.0,
	}

	result3 := ElementSerializer(obj3)
	assert.EqualValues(t, 3, len(result3))
	fmt.Println("LabMatrix", result3)
}
