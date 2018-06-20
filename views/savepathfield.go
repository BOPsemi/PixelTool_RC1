package views

import (
	"github.com/therecipe/qt/widgets"
)

/*
SavePathField :save path field
*/
type SavePathField struct {
	textLabel           *widgets.QLabel    // label
	textFieldForPath    *widgets.QLineEdit // text field for path
	textFieldForDirName *widgets.QLineEdit // text field for directory name
	textLabelForPath    *widgets.QLabel    // label for path

	DirectoryPath string
	DirectoryName string

	Cell *widgets.QWidget
}

/*
NewSavePathField :initializer of save path field
*/
func NewSavePathField(label string) *SavePathField {
	obj := new(SavePathField)

	// initialize widgets
	obj.Cell = widgets.NewQWidget(nil, 0)
	obj.textLabel = widgets.NewQLabel2(label, obj.Cell, 0)
	obj.textFieldForPath = widgets.NewQLineEdit(obj.Cell)
	obj.textFieldForDirName = widgets.NewQLineEdit(obj.Cell)
	obj.textLabelForPath = widgets.NewQLabel(obj.Cell, 0)

	// placeholder txt setup
	obj.textFieldForPath.SetPlaceholderText("Save Path")
	obj.textFieldForDirName.SetPlaceholderText("Dir Name")

	// layout
	layout := widgets.NewQGridLayout2()
	layout.AddWidget(obj.textLabel, 0, 0, 0)
	layout.AddWidget(obj.textFieldForPath, 0, 1, 0)
	layout.AddWidget(obj.textFieldForDirName, 0, 2, 0)
	layout.AddWidget3(obj.textLabelForPath, 1, 0, 1, 3, 0)

	// layout set
	obj.Cell.SetLayout(layout)

	// action connection
	obj.textFieldForPath.ConnectTextChanged(func(text string) {
		obj.DirectoryPath = text
		obj.textLabelForPath.SetText(obj.DirectoryPath)
	})

	obj.textFieldForDirName.ConnectTextChanged(func(text string) {
		obj.DirectoryName = text
		obj.textLabelForPath.SetText(obj.DirectoryPath + obj.DirectoryName)
	})

	return obj
}
