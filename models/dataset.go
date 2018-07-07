package models

/*
DataSet	:definition of deltaE data model
*/
type DataSet struct {
	DeltaE    []float64 // each patch deltaE data
	DivDeltaE []float64
	DeltaEAve float64   // averatege of deltaE
	Elm       []float64 // linear matrix elements
}

/*
WorstPatchNumber : rerutn worst patch number
*/
func (ds *DataSet) WorstPatchNumber() int {
	// check array size
	if len(ds.DeltaE) == 0 {
		return -1
	}

	// find worst data set
	number := 0
	worstVal := 0.0

	// find max
	for index, data := range ds.DeltaE {
		if data > worstVal {
			number = index
			worstVal = ds.DeltaE[number]
		}
	}

	// return
	return number

}
