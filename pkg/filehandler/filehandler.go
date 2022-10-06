package filehandler

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

// It opens the gzip file and uses the Reader to load the content
func GzipFileReader(gzFilePath string) *gzip.Reader {
	gzipFile, err := os.Open(gzFilePath)

	if err != nil {
		log.Fatal(err)
	}

	gzipReader, err := gzip.NewReader(gzipFile)

	if err != nil {
		log.Fatal(err)
	}

	gzipReader.Close()

	log.Print("Content opened of: " + gzFilePath)
	return gzipReader
}

// Creates an empty file if not exists and copy the content of the gzip.Reader to this one
func CopyGzipContentToFile(gzipReader io.Reader, destinationPath string) {
	fileWriter, err := os.Create(destinationPath)

	if err != nil {
		log.Fatal(err)
	}

	defer fileWriter.Close()

	_, err = io.Copy(fileWriter, gzipReader)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Content copied to: " + destinationPath)
}

// Get the file name without the last extension
func GetFileNameWithoutLastExtension(fileName string) string {
	lastExtension := filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(lastExtension)]
}

func JoinDirAndPath(dirName string, fileName string) string {
	return path.Join(filepath.Dir(dirName), fileName)
}

// Create empty file if not exists with JSON structure
func CreateEmptyJsonFileIfNotExists(fileName string) {
	_, err := os.Stat(fileName)

	if err != nil {
		if os.IsNotExist(err) {
			ioutil.WriteFile(fileName, []byte("{}"), 0)
		}
	}
}

// Get the content of a file and convert it from bytes to string
func GetFileContent(fileName string) string {
	fileContent, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer fileContent.Close()
	fileContentBytes, _ := ioutil.ReadAll(fileContent)
	return string(fileContentBytes)
}

func WriteFileContent(fileName string, data []byte) {
	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
