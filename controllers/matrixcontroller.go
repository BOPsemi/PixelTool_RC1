package controllers

import (
	"gonum.org/v1/gonum/mat"
)

/*
MatrixController :control matrix calculation
*/
type MatrixController interface {
	EvalLinearMatrix(elm []float64, grgbrb []float64) []float64
}

// object
type matrixController struct {
}

/*
NewMatrixController :initializer
*/
func NewMatrixController() MatrixController {
	obj := new(matrixController)

	return obj

}

/*
EvalLinearMatrix(elm []float64, rgb[]float64) []float64
*/
func (ma *matrixController) EvalLinearMatrix(elm []float64, grgbrb []float64) []float64 {
	result := make([]float64, 3)

	if (len(elm) == 6) && (len(grgbrb) == 4) {
		// result
		resultM := mat.NewDense(3, 1, make([]float64, 3))

		// matrix
		elmM := mat.NewDense(3, 3, ma.linmatGain(elm))
		rgbM := mat.NewDense(3, 1, ma.grgbrbToRGB(grgbrb))

		resultM.Mul(elmM, rgbM)

		// update
		result[0] = resultM.At(0, 0) // red
		result[1] = resultM.At(1, 0) // green
		result[2] = resultM.At(2, 0) // blue
	}

	return result
}

// GrGbRB -> RGB
func (ma *matrixController) grgbrbToRGB(grgbrb []float64) []float64 {
	/*
		grgbrb[0]	;Gr
		grgbrb[1]	;Gb
		grgbrb[2]	;r
		grgbrb[3]	;b
	*/

	/*
		Returns shold be this order
		R, G, B
	*/

	return []float64{grgbrb[2], (grgbrb[0] + grgbrb[1]) / 2.0, grgbrb[3]}
}

//
func (ma *matrixController) linmatGain(elm []float64) []float64 {
	rgain := 1.0 + elm[0] + elm[1] // red gain
	ggain := 1.0 + elm[2] + elm[3] // green gain
	bgain := 1.0 + elm[4] + elm[5] // blue gain

	// normarize all elements by each color gain
	rearrangedElm := []float64{
		1.0, -elm[0] / rgain, -elm[1] / rgain,
		-elm[2] / ggain, 1.0, -elm[3] / ggain,
		-elm[4] / bgain, -elm[5] / bgain, 1.0,
	}

	return rearrangedElm
}
