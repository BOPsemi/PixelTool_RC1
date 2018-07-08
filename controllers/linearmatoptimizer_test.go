package controllers

import (
	"PixelTool_RC1/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	linarMat = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/linearmatrix_elm.csv"
	dataPath = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	devQE    = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/device_QE.csv"
	ill      = models.D65
	gamma    = 0.45

	refCCPath = "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_color_patch.csv"
)

func Test_NewLinearMatrixOptimizer(t *testing.T) {
	obj := NewLinearMatrixOptimizer()

	assert.NotNil(t, obj)
}

func Test_SetEnv(t *testing.T) {
	obj := NewLinearMatrixOptimizer()
	assert.True(t, obj.SetEnv(linarMat, dataPath, devQE, ill, gamma))

}

func Test_SetRefColorCode(t *testing.T) {
	obj := NewLinearMatrixOptimizer()
	assert.True(t, obj.SetRefColorCode(refCCPath))

}

func Test_Run(t *testing.T) {
	obj := NewLinearMatrixOptimizer()

	// set enviroment for first run
	obj.SetEnv(linarMat, dataPath, devQE, ill, gamma)

	// set Ref CC
	obj.SetRefColorCode(refCCPath)

	// make linear matrix
	//elm := []float64{0.2201, 0.005, 0.0432, 0.0926, 0.00015, 0.398}
	elm := make([]float64, 0)

	// run
	obj.Run(elm, 3.0, 10, 1.0, 5)

	// result
	fmt.Println(obj.OptimizedLinearMatrix())
	deltaE, deltaEAve := obj.FinalDeltaEInfo()
	fmt.Println(deltaEAve, deltaE)

}
