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
package shipper

import (
	"beatshipper/configs"
	"beatshipper/pkg/filehandler"
	"beatshipper/pkg/registry"
	"log"
	"os"
	"time"

	client "github.com/elastic/go-lumber/client/v2"
)

func setup() (configs.Configuration, *registry.Registry) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := configs.Configuration{}
	c.GetConfiguration()

	filehandler.CreateEmptyJsonFileIfNotExists(c.Registry)
	r := &registry.Registry{}

	// We need to re-read the content of the registry each time we wanna check if
	// there are new files to process
	registryFileContent := filehandler.GetFileContent(c.Registry)
	registry.ConvertJSONToRegistry([]byte(registryFileContent), r)

	return c, r
}

func Run() {
	config, register := setup()

	// Globbing all paths according to the pattern
	files := filehandler.FoundFilesByPaths(config.Paths)

	if len(files) < 1 {
		log.Print("There aren't files to process")
		os.Exit(0)
	}

	filesToProcess := getFilesNotProcessed(files, *register)

	if len(filesToProcess) < 1 {
		log.Print("There aren't files to process")
	} else {
		SendBatch(filesToProcess, config)
		registry.StoreFilesIntoRegistry(filesToProcess, *register, config.Registry)
	}
}

// Creates the conexion
// Get the content of files and create a batch with event interface
// Send the batch to the server
// Closes the conexion
func SendBatch(files []string, config configs.Configuration) error {
	conn, err := client.SyncDial(config.Host + ":" + config.Port)

	if err != nil {
		log.Fatal(err, ": Check if the host is listening")
	}

	log.Print("Conexion successful with: ", config.Host+":"+config.Port)

	var batch []interface{}

	for _, filePath := range files {
		log.Print("Processing: " + filePath)
		batch = append(batch, makeEvent(filehandler.GetFileContentByExtension(filePath)))
	}

	chunkedBatch := chunkInterfaceSlice(batch, 5)

	for _, batchSlices := range chunkedBatch {

		_, err = conn.Send(batchSlices)

		if err != nil {
			log.Fatal(err)
		}

		log.Print("Sending batch of data...")
	}

	err = conn.Close()

	log.Print("Conexion closed")

	return err
}

func getFilesNotProcessed(files []string, r registry.Registry) []string {
	filesNotProcessed := []string{}
	for _, filePath := range files {
		if !r.IsFileInRegistry(filePath) {
			filesNotProcessed = append(filesNotProcessed, filePath)
		}
	}
	return filesNotProcessed
}

// Chunk a slice into uniform slices according to the size
func chunkInterfaceSlice(slice []interface{}, size int) [][]interface{} {
	log.Print("Chunk slice into: ", size, " slices")

	var chunks [][]interface{}

	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// Create a valid event that can be received by logstash
func makeEvent(message string) interface{} {
	return map[string]interface{}{
		"@timestamp": time.Now(),
		"type":       "filebeat",
		"message":    message,
	}
}
