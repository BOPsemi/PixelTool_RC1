package util

/*
MakeMacbethBoard :interface definition of MakeMacbethBoard
*/
type MakeMacbethBoard interface {
	SetDataSource(path string) bool
}

// definition of structure
type makeMacbethBoard struct {
	fileNames []string
}

/*
NewMakeMacbethBoard :initializer
*/
func NewMakeMacbethBoard() MakeMacbethBoard {
	obj := new(makeMacbethBoard)

	// initialize properties
	obj.fileNames = make([]string, 0)

	return obj
}

/*
SetDataSource : setter of data source
	in	;path string
	out	;bool
*/
func (mm *makeMacbethBoard) SetDataSource(path string) bool {
	status := false

	return status
}
