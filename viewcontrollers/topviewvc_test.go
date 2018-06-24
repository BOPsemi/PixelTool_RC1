package viewcontrollers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTopViewViewController(t *testing.T) {
	obj := NewTopViewViewController()

	assert.NotNil(t, obj)
}

func Test_EvaluateDeltaE(t *testing.T) {
	obj := NewTopViewViewController()

	refPath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/std_patch/std_24_ColorChart.csv"
	devPath := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/dev_patch/dev_24_ColorChart.csv"

	kvalues := []float64{1.0, 1.0, 1.0}

	results, status := obj.EvaluateDeltaE(refPath, devPath, kvalues)
	assert.True(t, status)
	assert.EqualValues(t, 24, len(results))

	status = obj.SaveDeltaEResultData()
	assert.True(t, status)

}
