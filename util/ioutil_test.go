package util

import "testing"
import "github.com/stretchr/testify/assert"

func Test_NewIOUtil(t *testing.T) {
	obj := NewIOUtil()

	assert.NotNil(t, obj)
}

func Test_WriteCSVFile(t *testing.T) {
	obj := NewIOUtil()

	path := "/Users/kazufumiwatanabe/go/src/PixelTool/data/"
	filename := "writeCSVtest"

	line1 := []string{"1", "hoge", "hige"}
	line2 := []string{"2", "red", "blue"}
	data := [][]string{line1, line2}

	obj.WriteCSVFile(path, filename, data)

}
