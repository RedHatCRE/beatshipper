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
package main

import (
	"gz-beat-shipper/configs"
	"gz-beat-shipper/pkg/filehandler"
	"gz-beat-shipper/pkg/registry"
	"log"
	"path/filepath"
	"time"

	client "github.com/elastic/go-lumber/client/v2"
)

func setup() (configs.Configuration, *registry.Registry) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	c := configs.Configuration{}
	c.GetConfiguration()

	filehandler.CreateEmptyJsonFileIfNotExists(c.Registry)
	r := &registry.Registry{}

	return c, r
}

func main() {

	config, r := setup()

	recheckTime := config.GetRecheckDuration()

	conn, err := client.SyncDial(config.Host + ":" + config.Port)

	if err != nil {
		log.Fatal(err)
	}

	for {
		// We need to re-read the content of the registry each time we wanna check if
		// there are new files to process
		registryFileContent := filehandler.GetFileContent(config.Registry)
		registry.ConvertJSONToRegistry([]byte(registryFileContent), r)

		// Globbing all paths according to the pattern
		gzFilesPaths, err := filepath.Glob(config.Path)

		if err != nil {
			log.Fatal(err)
		}

		if len(gzFilesPaths) < 1 {
			log.Print("There aren't gz files to process for " + config.Path)
			time.Sleep(recheckTime)
			continue
		}

		gzFilesToProcess := getGzNotProcessed(gzFilesPaths, *r)

		if len(gzFilesToProcess) < 1 {
			log.Print("There aren't gz files to process for " + config.Path)
			time.Sleep(recheckTime)
			continue
		}

		_, err = conn.Send(getBatchToSend(gzFilesToProcess))

		if err != nil {
			log.Fatal(err)
		}

		registry.StoreFilesIntoRegistry(gzFilesToProcess, *r, config.Registry)

		time.Sleep(recheckTime)
	}

}

func getGzNotProcessed(gzFiles []string, r registry.Registry) []string {
	gzFilesNotProcessed := []string{}
	for _, gzFilePath := range gzFiles {
		if !r.IsFileInRegistry(gzFilePath) {
			gzFilesNotProcessed = append(gzFilesNotProcessed, gzFilePath)
		}
	}
	return gzFilesNotProcessed
}

// Unzip the Gz in the same directory where they are
func getBatchToSend(gzFiles []string) []interface{} {
	batch := make([]interface{}, 1)
	i := 0
	for _, gzFilePath := range gzFiles {
		log.Print("Processing: " + gzFilePath)
		batch[i] = makeEvent(filehandler.GetGzFileContent(gzFilePath))
	}
	return batch
}

// Create a valid event that can be received by logstash
func makeEvent(message string) interface{} {
	return map[string]interface{}{
		"@timestamp": time.Now(),
		"type":       "filebeat",
		"message":    message,
	}
}
