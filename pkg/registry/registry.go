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
package registry

import (
	"beatshipper/pkg/filehandler"
	"encoding/json"
	"log"
	"time"
)

type Registry struct {
	ParsedFiles []ParsedFiles
}

type ParsedFiles struct {
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

// Append File struct to Registry Struct with ParsedFiles slice
func (r *Registry) AppendFileToRegistry(f string) {
	r.ParsedFiles = append(r.ParsedFiles, ParsedFiles{
		Name: f,
		Date: time.Now(),
	})
}

// Convert Registry Struct to a JSON string
func (r *Registry) JSON() string {
	jsonRegistry, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonRegistry)
}

// Convert a JSON from bytes to the Registry struct format
func ConvertJSONToRegistry(jsonRegistry []byte, registry *Registry) {
	if err := json.Unmarshal(jsonRegistry, registry); err != nil {
		log.Fatal(err)
	}
}

// Check if a file in a string way exists in the Registry struct format
func (r *Registry) IsFileInRegistry(f string) bool {
	for _, parsedFile := range r.ParsedFiles {
		if parsedFile.Name == f {
			return true
		}
	}
	return false
}

/*
	Store all the provided files in the Registry struct
	Convert the struct to JSON byte data
	Write the content in the registry file
*/
func StoreFilesIntoRegistry(fileNames []string, r Registry, registryFileLocation string) {
	for _, f := range fileNames {
		r.AppendFileToRegistry(f)
		log.Print("Added to registry: ", f)
	}
	j := r.JSON()
	filehandler.WriteFileContent(registryFileLocation, []byte(j))
}
