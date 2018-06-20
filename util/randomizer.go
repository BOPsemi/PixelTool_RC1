package util

import (
	"math/rand"
	"time"
)

/*
Randomizer :interface of randomizer
*/
type Randomizer interface {
	Run(list []int, iteration int) []int
}

// object
type randomizer struct {
}

/*
NewRandomizer : initializer of randamizer
*/
func NewRandomizer() Randomizer {
	obj := new(randomizer)

	return obj
}

/*
Run(iteration int) []int
*/
func (ra *randomizer) Run(list []int, iteration int) []int {
	// initialize result slice, just fill zero to result
	result := make([]int, 0)
	for i := 0; i < len(list); i++ {
		result = append(result, 0)
	}

	// initialize random seed
	rand.Seed(time.Now().UnixNano())

	// sum of list
	sumOflist := func(datalist []int) int {
		sum := 0
		for _, data := range datalist {
			sum += data
		}

		return sum
	}

	// count index
	countIndex := func(datalist []int) int {
		retIndex := -1

		// random number
		rnumber := rand.Intn(sumOflist(datalist) + 1)

		// lottery
		for index, value := range datalist {
			if value >= rnumber {
				retIndex = index
				break
			}
			rnumber -= value
		}

		return retIndex
	}

	// start randomize
	for i := 0; i < iteration; i++ {
		index := countIndex(list)
		result[index]++
	}

	return result
}
