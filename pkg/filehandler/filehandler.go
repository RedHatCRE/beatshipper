// Copyright 2022 Red Hat, Inc.
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.
package filehandler

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

// It opens the gzip file and uses the Reader to load the content
func GzipFileReader(gzFilePath string) *gzip.Reader {

	gzipFile, err := os.Open(filepath.Clean(gzFilePath))

	if err != nil {
		log.Fatal(err)
	}

	gzipReader, err := gzip.NewReader(gzipFile)

	if err != nil {
		log.Fatal(err)
	}

	err = gzipReader.Close()

	if err != nil {
		log.Print("Error closing filehandle: ", err)
	}

	log.Print("Content opened of: " + gzFilePath)
	return gzipReader
}

// Creates an empty file if not exists and copy the content of the gzip.Reader to this one
func CopyGzipContentToFile(gzipReader io.Reader, destinationPath string) {
	fileWriter, err := os.Create(filepath.Clean(destinationPath))

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := fileWriter.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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

	if errors.Is(err, os.ErrNotExist) {
		err = ioutil.WriteFile(fileName, []byte("{}"), 0600)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

}

func GetFileContentByExtension(fileName string) string {
	extension := filepath.Ext(fileName)
	if extension == ".gz" {
		return GetGzFileContent(fileName)
	} else {
		return GetFileContent(fileName)
	}
}

// Get the content of a file and convert it from bytes to string
func GetFileContent(fileName string) string {
	fileContent, err := os.Open(filepath.Clean(fileName))

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := fileContent.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fileContentBytes, _ := ioutil.ReadAll(fileContent)
	return string(fileContentBytes)
}

func WriteFileContent(fileName string, data []byte) {
	err := ioutil.WriteFile(fileName, data, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func GetGzFileContent(fileName string) string {
	fh, err := os.Open(filepath.Clean(fileName))

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := fh.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fhGz, err := gzip.NewReader(fh)
	if err != nil {
		log.Fatal(err)
	}
	defer fhGz.Close()

	contentBytes, err := ioutil.ReadAll(fhGz)
	if err != nil {
		log.Fatal(err)
	}

	return string(contentBytes)
}

// Find files based in a list of paths  using glob
func FoundFilesByPaths(paths []string) []string {
	var fileNames []string
	for _, path := range paths {
		names, err := filepath.Glob(path)

		if err != nil {
			log.Fatal(err)
		}

		fileNames = append(fileNames, names...)
	}

	return fileNames
}
