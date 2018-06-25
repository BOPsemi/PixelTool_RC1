package util

/*
LinearMatrixOptizer :optimzer
*/
type LinearMatrixOptizer interface {
	SetDataSet(elm []float64, deltaE []float64)
	Run()
}

// struct
type linearMatrixOptizer struct {
	defaultElm    []float64
	defaultDeltaE []float64
	defaultEave   float64
}

/*
NewLinearMatrixOptizer :initializer
*/
func NewLinearMatrixOptizer() LinearMatrixOptizer {
	obj := new(linearMatrixOptizer)

	return obj
}

/*
SetDataSet
	in	:elm []float64, deltaE []float64
	out	:
*/
func (lo *linearMatrixOptizer) SetDataSet(elm []float64, deltaE []float64) {
	// initialize data
	lo.defaultElm = elm
	lo.defaultDeltaE = deltaE

	// calculate deltaE average
	ave := 0.0
	for _, data := range lo.defaultDeltaE {
		ave += data
	}
	lo.defaultEave = ave / float64(len(lo.defaultDeltaE))
}

/*
Run
*/
func (lo *linearMatrixOptizer) Run() {

}
