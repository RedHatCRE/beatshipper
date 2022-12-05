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
	"errors"
	"log"
	"os"
	"syscall"
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
	defer timeTrack(time.Now(), "Beatshipper")

	config, memoryRegistry := setup()
	registryFileLocation := config.Registry

	// Globbing all paths according to the pattern
	files := filehandler.FoundFilesByPaths(config.Paths)

	if len(files) < 1 {
		log.Fatal("There aren't files to process")
	}

	filesToProcess := getFilesNotProcessed(files, *memoryRegistry)

	if len(filesToProcess) < 1 {
		log.Fatal("There aren't files to process")
	}

	processData(filesToProcess, config, memoryRegistry, registryFileLocation)
}

// Creates the conexion
// Get the content of files and create a batch with event interface
// Send the batch to the server
// Closes the conexion
func processData(filesToProcess []string, config configs.Configuration, memoryRegistry *registry.Registry, registryFileLocation string) {
	conn, err := client.SyncDial(config.Host + ":" + config.Port)

	if err != nil {
		log.Fatal(err, ": Check if the host is listening")
	}

	log.Print("Conexion successful with: ", config.Host+":"+config.Port)

	hostName, err := os.Hostname()

	if err != nil {
		log.Fatal(err)
	}

	var batch []interface{}

	for _, filePath := range filesToProcess {
		data := map[string]interface{}{
			"message": filehandler.GetFileContentByExtension(filePath),
			"host": map[string]interface{}{
				"name": hostName,
			},
			"log": map[string]interface{}{
				"file": map[string]interface{}{
					"path": filePath,
				},
			},
			"@timestamp": time.Now(),
			"log_source": config.LogSource,
			"type":       "beatshipper",
		}

		batch = append(batch, data)
	}

	// Chunk the amount of batches into a slice to send them in different groups and not
	// load the destination
	chunkedBatch := chunkInterfaceSlice(batch, 30)

	var processedFiles []string

	for _, batchSlices := range chunkedBatch {

		_, err = conn.Send(batchSlices)

		if err != nil {
			// If the destination service configuration is wrong but the service is still up (it was able to connect)
			// it won't be able to process the information so we should handle the error that is more common
			if errors.Is(err, syscall.EPIPE) {
				log.Printf("EPIPE error. Check the destination configuration [%s]", err)
			} else {
				log.Print(err)
			}
			continue
		}

		log.Print("Sending batch of data...")

		// We need to add into the registry just the files that have been sent successfully
		appendFilesInSliceOfBatch(&processedFiles, batchSlices)
	}

	err = conn.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Conexion closed")

	if len(processedFiles) > 0 {
		registry.StoreFilesIntoRegistry(processedFiles, *memoryRegistry, registryFileLocation)
	}
}

// We need to access to the files of a nested map[string]interface{} data, for this we should use a type of reflection
// to access to every key of the nested map
func appendFilesInSliceOfBatch(files *[]string, batch []interface{}) {
	for i := 0; i < len(batch); i += 1 {
		*files = append(*files, batch[i].(map[string]interface{})["log"].(map[string]interface{})["file"].(map[string]interface{})["path"].(string))
	}
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

func timeTrack(start time.Time, name string) {
	log.Printf("%s took %s", name, time.Since(start))
}
