package views

import "github.com/therecipe/qt/widgets"

/*
FileOutputField :File output field
*/
type FileOutputField struct {
	textLabel            *widgets.QLabel    // label
	textFieldForPath     *widgets.QLineEdit // text field for path
	textFieldForFileName *widgets.QLineEdit // text field for directory name
	textLabelForExt      *widgets.QLabel    // label for path
	textLabelForFull     *widgets.QLabel    // label for full path

	DirectoryPath string
	FileName      string

	Cell *widgets.QWidget
}

/*
NewFileOutputField :initializer of file output field
*/
func NewFileOutputField(label string) *FileOutputField {
	obj := new(FileOutputField)

	// initialize widgets
	obj.Cell = widgets.NewQWidget(nil, 0)
	obj.textLabel = widgets.NewQLabel2(label, obj.Cell, 0)
	obj.textFieldForPath = widgets.NewQLineEdit(obj.Cell)
	obj.textFieldForFileName = widgets.NewQLineEdit(obj.Cell)
	obj.textLabelForExt = widgets.NewQLabel(obj.Cell, 0)
	obj.textLabelForFull = widgets.NewQLabel(obj.Cell, 0)

	// placeholder txt setup
	obj.textFieldForPath.SetPlaceholderText("Save Path")
	obj.textFieldForFileName.SetPlaceholderText("File Name")
	obj.textLabelForExt.SetText(".csv")

	// layout
	layout := widgets.NewQGridLayout2()
	layout.AddWidget(obj.textLabel, 0, 0, 0)
	layout.AddWidget(obj.textFieldForPath, 0, 1, 0)
	layout.AddWidget(obj.textFieldForFileName, 0, 2, 0)
	layout.AddWidget(obj.textLabelForExt, 0, 3, 0)
	layout.AddWidget3(obj.textLabelForFull, 1, 0, 1, 3, 0)

	// layout set
	obj.Cell.SetLayout(layout)

	// action connection
	obj.textFieldForPath.ConnectTextChanged(func(text string) {
		obj.DirectoryPath = text
		obj.textLabelForFull.SetText(obj.DirectoryPath)
	})

	obj.textFieldForFileName.ConnectTextChanged(func(text string) {
		obj.FileName = text
		obj.textLabelForFull.SetText(obj.DirectoryPath + obj.FileName + ".csv")
	})

	return obj
}
