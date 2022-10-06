package filehandler

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Check GNU zip decompress func works
func TestGzipFileReader(t *testing.T) {
	gzTestFileName := "test.txt.gz"
	gzipReader := GzipFileReader(gzTestFileName)
	text, _ := ioutil.ReadAll(gzipReader)
	assert := assert.New(t)
	assert.Equal("This is the content of the compressed file\n", string(text))
}

// Check gzipReader copy the content to the destination file
// It should create it and remove it
func TestCopyGzipContentToFile(t *testing.T) {
	gzTestFileName := "test.txt.gz"
	destinationFile := "test.txt"
	gzipReader := GzipFileReader(gzTestFileName)
	CopyGzipContentToFile(gzipReader, destinationFile)

	assert := assert.New(t)
	assert.True(existFile(destinationFile))

	_ = os.Remove(destinationFile)
}

// Test if remove the last extension in a right way
func TestGetFileNameWithoutLastExtension(t *testing.T) {
	gzFileName := "test.txt.gz"
	txtFileName := GetFileNameWithoutLastExtension(gzFileName)
	assert := assert.New(t)
	assert.Equal("test.txt", txtFileName)
	fileName := GetFileNameWithoutLastExtension(txtFileName)
	assert.Equal("test", fileName)
}

// There's no a direct way to check if file exists in Go
// So we should stat it with a call to the system
func existFile(fileName string) bool {
	_, err := os.Stat(fileName)

	// check if error is "file not exists"
	if os.IsNotExist(err) {
		return false
	}
	return true
}
