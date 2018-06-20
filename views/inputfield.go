package views

import (
	"github.com/therecipe/qt/widgets"
)

/*
InputField :define input field widget
*/
type InputField struct {
	textLabel *widgets.QLabel      // text label
	textField *widgets.QLineEdit   // text filed for path phrase
	button    *widgets.QPushButton // clear button

	InputedText string // text data in the text field

	Cell *widgets.QWidget // input field widget
}

/*
NewInputField : initalizer of input field structure
*/
func NewInputField(label, placeholder string) *InputField {
	obj := new(InputField)

	// initialize widgets
	obj.Cell = widgets.NewQWidget(nil, 0)
	obj.textLabel = widgets.NewQLabel2(label, obj.Cell, 0)
	obj.textField = widgets.NewQLineEdit(obj.Cell)
	obj.textField.SetPlaceholderText(placeholder)
	obj.button = widgets.NewQPushButton2("Clear", obj.Cell)

	// layout
	layout := widgets.NewQHBoxLayout()
	layout.AddWidget(obj.textLabel, 0, 0)
	layout.AddWidget(obj.textField, 0, 0)
	layout.AddWidget(obj.button, 0, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// connect action
	obj.button.ConnectClicked(func(checked bool) {
		obj.textField.Clear()
		obj.textField.Repaint()
	})

	obj.textField.ConnectTextChanged(func(text string) {
		obj.InputedText = obj.textField.Text()
	})

	return obj
}
