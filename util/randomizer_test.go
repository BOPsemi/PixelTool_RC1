package util

import "testing"
import "github.com/stretchr/testify/assert"
import "fmt"

func Test_NewRandomizer(t *testing.T) {
	obj := NewRandomizer()
	assert.NotNil(t, obj)

}

func Test_Run(t *testing.T) {
	obj := NewRandomizer()

	list := []int{80, 10, 8, 1, 1}
	result := obj.Run(list, 1000)

	fmt.Println(result)
}
