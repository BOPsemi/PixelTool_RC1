package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDeltaLabCalculator(t *testing.T) {
	obj := NewDeltaLabCalculator()
	assert.NotNil(t, obj)
}

func Test_DeltaLab(t *testing.T) {
	obj := NewDeltaLabCalculator()

	ref := []float64{100.0, 0.0, 0.0}
	comp := []float64{50.0, -10.0, 1.0}
	kvalues := []float64{1.0, 1.0, 1.0}

	distance := obj.DeltaLab(ref, comp, kvalues)

	fmt.Println(distance)
}
