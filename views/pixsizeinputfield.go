package views

import (
	"strconv"

	"github.com/therecipe/qt/widgets"
)

/*
PixelSizeInputField :pixel and patch size input field
*/
type PixelSizeInputField struct {
	texLabel      *widgets.QLabel // title
	textLabelForH *widgets.QLabel // horizontal
	textLabelForV *widgets.QLabel // vertical

	textFieldForH *widgets.QLineEdit // field for horizontal
	textFieldForV *widgets.QLineEdit // field for vertical

	Cell           *widgets.QWidget
	HorizontalSize int
	VerticalSize   int
}

/*
NewPixelSizeInputField :initializer
*/
func NewPixelSizeInputField(label string, horizontal, vertical int) *PixelSizeInputField {
	obj := new(PixelSizeInputField)
	obj.HorizontalSize = horizontal
	obj.VerticalSize = vertical

	// initialize widgets
	obj.Cell = widgets.NewQWidget(nil, 0)
	obj.texLabel = widgets.NewQLabel2(label, obj.Cell, 0)
	obj.textLabelForH = widgets.NewQLabel2("H:", obj.Cell, 0)
	obj.textLabelForV = widgets.NewQLabel2("V:", obj.Cell, 0)

	obj.textFieldForH = widgets.NewQLineEdit2(strconv.Itoa(horizontal), obj.Cell)
	obj.textFieldForV = widgets.NewQLineEdit2(strconv.Itoa(vertical), obj.Cell)

	// layout
	layout := widgets.NewQHBoxLayout()
	layout.AddWidget(obj.texLabel, 0, 0)
	layout.AddWidget(obj.textLabelForH, 0, 0)
	layout.AddWidget(obj.textFieldForH, 0, 0)
	layout.AddWidget(obj.textLabelForV, 0, 0)
	layout.AddWidget(obj.textFieldForV, 0, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// action connection
	obj.textFieldForH.ConnectTextChanged(func(text string) {
		obj.HorizontalSize = obj.stringToIntConverter(text)

		//fmt.Println(obj.HorizontalSize)
	})

	obj.textFieldForV.ConnectTextChanged(func(text string) {
		obj.VerticalSize = obj.stringToIntConverter(text)

		//fmt.Println(obj.VerticalSize)
	})

	return obj
}

// string to int value converter
func (pf *PixelSizeInputField) stringToIntConverter(str string) int {
	val, err := strconv.Atoi(str)

	// check error
	if err != nil {
		return 100
	}

	// check negative value
	if val < 0 {
		return 100
	}

	// return value
	return val
}
