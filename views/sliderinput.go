package views

import (
	"strconv"

	"github.com/therecipe/qt/widgets"
)

/*
SliderInput :define slider
*/
type SliderInput struct {
	slider    *widgets.QSlider   // slider
	textLabel *widgets.QLabel    // text label
	textField *widgets.QLineEdit // text fiedl for direct input

	Value float64

	Cell *widgets.QWidget
}

/*
NewSliderInput :initializer of slider input
*/
func NewSliderInput(label string, initValue float64) *SliderInput {
	obj := new(SliderInput)

	//initialize widgets
	obj.Cell = widgets.NewQWidget(nil, 0)
	obj.slider = widgets.NewQSlider2(1, obj.Cell)
	obj.textLabel = widgets.NewQLabel2(label, obj.Cell, 0)
	obj.textField = widgets.NewQLineEdit(obj.Cell)

	layout := widgets.NewQHBoxLayout()
	layout.AddWidget(obj.textLabel, 0, 0)
	layout.AddWidget(obj.slider, 0, 0)
	layout.AddWidget(obj.textField, 0, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// set initial value
	obj.slider.SetValue(int(100.0 * initValue))
	str := strconv.FormatFloat(initValue, 'f', 2, 64)
	obj.textField.SetText(str)

	obj.Value = initValue

	// action connection
	obj.slider.ConnectSliderMoved(func(value int) {
		normedValue := float64(value) / 100.0
		str := strconv.FormatFloat(normedValue, 'f', 2, 64)

		obj.textField.SetText(str)

		// upload value
		obj.Value = normedValue

	})

	obj.textField.ConnectTextChanged(func(text string) {
		var normedValue int
		value, _ := strconv.ParseFloat(text, 64)
		if value > 1.0 {
			normedValue = 100
		} else {
			normedValue = int(value * 100)
		}

		obj.slider.SetValue(normedValue)

		// upload value
		obj.Value = float64(normedValue) / 100.0
	})

	return obj
}
