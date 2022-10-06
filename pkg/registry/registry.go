package registry

import (
	"encoding/json"
	"gz-beat-shipper/pkg/filehandler"
	"log"
)

type Registry struct {
	ParsedFiles []File `json:"ParsedFile"`
}

type File struct {
	Name string `json:"name"`
}

// Append File struct to Registry Struct with ParsedFiles slice
func (r *Registry) AppendFileToRegistry(f File) {
	r.ParsedFiles = append(r.ParsedFiles, f)
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
		f := File{
			Name: f,
		}
		r.AppendFileToRegistry(f)
		log.Print("Added to registry: ", f)
	}
	j := r.JSON()
	filehandler.WriteFileContent(registryFileLocation, []byte(j))
}
