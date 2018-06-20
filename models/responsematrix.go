package models

import (
	"reflect"
)

/*
RGBMatrix :define RGB matrix
*/
type RGBMatrix struct {
	R float64
	G float64
	B float64
}

/*
XYZMatrix :define XYZ matrix
*/
type XYZMatrix struct {
	X float64
	Y float64
	Z float64
}

/*
LabMatrix :define Lab matrix
*/
type LabMatrix struct {
	L float64
	A float64
	B float64
}

/*
ElementSerializer :serializer
*/
func ElementSerializer(obj interface{}) []float64 {

	objType := reflect.TypeOf(obj).String()

	switch objType {
	case "*models.RGBMatrix":
		// RGB matrix case
		if object, ok := obj.(*RGBMatrix); ok {
			return []float64{object.R, object.G, object.B}
		}
		return []float64{}

	case "*models.XYZMatrix":
		// XYZ matrix case
		if object, ok := obj.(*XYZMatrix); ok {
			return []float64{object.X, object.Y, object.Z}
		}
		return []float64{}

	case "*models.LabMatrix":
		// Lab matrix case
		if object, ok := obj.(*LabMatrix); ok {
			return []float64{object.L, object.A, object.B}
		}
		return []float64{}

	default:
		return []float64{}
	}

}
