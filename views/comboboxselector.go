package views

import (
	"github.com/therecipe/qt/widgets"
)

/*
ComboBoxSelector :combo box selector
*/
type ComboBoxSelector struct {
	box       *widgets.QComboBox // combo box
	textLabel *widgets.QLabel    // label

	Cell         *widgets.QWidget
	SelectedItem string
}

/*
NewComboBoxSelector :initializer of combo box selector
*/
func NewComboBoxSelector(label string, list []string) *ComboBoxSelector {
	obj := new(ComboBoxSelector)
	obj.SelectedItem = list[0]

	// initialize widgets
	obj.Cell = widgets.NewQWidget(nil, 0)
	obj.box = widgets.NewQComboBox(obj.Cell)
	obj.box.AddItems(list)
	obj.textLabel = widgets.NewQLabel2(label, obj.Cell, 0)

	// layout
	layout := widgets.NewQHBoxLayout()
	layout.AddWidget(obj.textLabel, 0, 0)
	layout.AddWidget(obj.box, 0, 0)

	// apply layout
	obj.Cell.SetLayout(layout)

	// action connection
	obj.box.ConnectCurrentIndexChanged(func(index int) {
		obj.SelectedItem = list[index]

		//fmt.Println(obj.SelectedItem)
	})

	return obj
}
