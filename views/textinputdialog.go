package views

import "github.com/therecipe/qt/widgets"

/*
TextInputDialog :text input dialog
*/
type TextInputDialog struct {
	Cell *widgets.QInputDialog
}

/*
NewTextInputDialog : initalizer
*/
func NewTextInputDialog(title string, label string) *TextInputDialog {
	obj := new(TextInputDialog)

	obj.Cell = widgets.NewQInputDialog(nil, 0)
	obj.Cell.Resize2(400, 300)
	obj.Cell.SetWindowTitle(title)
	obj.Cell.SetLabelText(label)

	return obj
}
