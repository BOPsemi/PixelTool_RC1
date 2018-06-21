package viewcontrollers

/*
MainViewController : main view controller
*/
type MainViewController interface {
}

// mainViewController
type mainViewController struct {
}

/*
NewMainViewController : initializer
*/
func NewMainViewController() MainViewController {
	obj := new(mainViewController)

	return obj
}
