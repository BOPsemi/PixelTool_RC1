package controllers

import (
	"PixelTool_RC1/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewResponseController(t *testing.T) {
	obj := NewResponseController()

	assert.NotNil(t, obj)
}

func Test_ReadResponseData(t *testing.T) {
	obj := NewResponseController()

	// make file list for testing
	path := make(map[string]string, 0)
	path["DeviceQE"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/device_QE.csv"
	path["ColorChecker"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/Macbeth_Color_Checker.csv"
	path["D65"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_D65.csv"
	path["IllA"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_A.csv"

	status := obj.ReadResponseData(path)

	assert.True(t, status)

}

func Test_CalculateResponse(t *testing.T) {
	obj := NewResponseController()

	// make file list for testing
	path := make(map[string]string, 0)
	path["DeviceQE"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/device_QE.csv"
	path["ColorChecker"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/Macbeth_Color_Checker.csv"
	path["D65"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_D65.csv"
	path["IllA"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_A.csv"

	if obj.ReadResponseData(path) {
		status, responses := obj.CalculateChannelResponse(models.D65, 400, 700, 5, 22)
		if status {

			fmt.Println("-------------")
			for _, data := range responses {
				flag, result := obj.CalculateGammaCorrection(0.43, &data)

				if flag {
					fmt.Println(data.CheckerNumber, result.Gr, result.Gb, result.R, result.B)
				}
			}
		}
	}
}

func Test_CalculateGammaCorrection(t *testing.T) {
	obj := NewResponseController()

	// make file list for testing
	path := make(map[string]string, 0)
	path["DeviceQE"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/device_QE.csv"
	path["ColorChecker"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/Macbeth_Color_Checker.csv"
	path["D65"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_D65.csv"
	path["IllA"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_A.csv"

	if obj.ReadResponseData(path) {
		status, responses := obj.CalculateChannelResponse(models.D65, 400, 700, 5, 22)

		if status {

			for index, data := range responses {
				success, res := obj.CalculateGammaCorrection(0.43, &data)
				if success {
					fmt.Println(index+1, res)
				}
			}
		}
	}
}

func Test_CalculateLinearMatrix(t *testing.T) {

	obj := NewResponseController()

	// make file list for testing
	path := make(map[string]string, 0)
	path["DeviceQE"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/device_QE.csv"
	path["ColorChecker"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/Macbeth_Color_Checker.csv"
	path["D65"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_D65.csv"
	path["IllA"] = "/Users/kazufumiwatanabe/go/src/PixelTool/data/illumination_A.csv"

	if obj.ReadResponseData(path) {
		status, responses := obj.CalculateChannelResponse(models.D65, 400, 700, 5, 22)
		if status {
			// stocker of data which was correcd by gamma
			gammaCorrectedRes := make([]models.ChannelResponse, 0)

			for _, data := range responses {
				flag, result := obj.CalculateGammaCorrection(0.43, &data)
				if flag {
					gammaCorrectedRes = append(gammaCorrectedRes, *result)
				}
			}

			// calculate linear matrix
			elm := []float64{0.136, 0.061, 0.104, 0.063, 0.041, 0.057}
			results := make([][]float64, 0)

			for _, data := range gammaCorrectedRes {
				grgbrb := []float64{
					data.Gr,
					data.Gb,
					data.R,
					data.B,
				}

				result := obj.CalculateLinearMatrix(elm, grgbrb)
				results = append(results, result)
			}

			// calculate Red/Blue gain
			redGain, blueGain := obj.CalculateWhiteBalanceGain(results[21])
			fmt.Println(redGain, blueGain)

		}
	}

}
