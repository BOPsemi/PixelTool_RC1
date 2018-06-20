package util

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewDirectoryHandler(t *testing.T) {
	obj := NewDirectoryHandler()

	assert.NotNil(t, obj)
}

func Test_GetCurrentDirectoryPath(t *testing.T) {
	obj := NewDirectoryHandler()
	path := obj.GetCurrentDirectoryPath()

	assert.NotEmpty(t, path)
	log.Println(path)
}

func Test_GetFileListInDirectory(t *testing.T) {
	obj := NewDirectoryHandler()
	path := obj.GetCurrentDirectoryPath()

	_, files := obj.GetFileListInDirectory(path)
	//assert.EqualValues(t, 0, len(files))

	if len(files) > 0 {
		log.Println(files)

	} else {

		fmt.Println("no files")
	}
}

func Test_MakeDirectory(t *testing.T) {
	obj := NewDirectoryHandler()

	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1"
	assert.True(t, obj.MakeDirectory(path, "hoge"))
}

func Test_DirectoryAvailable(t *testing.T) {
	obj := NewDirectoryHandler()

	path := "/Users/kazufumiwatanabe/go/src/PixelTool_RC1/data/"
	//path := "/Users/kazufumiwatanabe/go/src/PixelTool/data/device_QE.csv"
	//path := "/Users/kazufumiwatanabe/go/src/PixelTool/data/doc"
	assert.False(t, obj.DirectoryAvailable(path))
}
