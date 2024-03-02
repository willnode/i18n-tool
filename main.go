package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v3"
)

func main() {

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	RealMain(currentDir)
}

func RealMain(currentDir string) error {

	// Define folder paths
	sourceFolder := "en"
	targetFolders, err := getListOfLanguage(currentDir)
	if err != nil {
		log.Fatalf("Failed to get list of target folders: %v", err)
		return err
	}

	// Construct full path for source folder
	sourceFolderPath := filepath.Join(currentDir, sourceFolder)

	// Iterate over target folders
	for _, targetFolder := range targetFolders {

		// Traverse folder A
		err := filepath.Walk(sourceFolderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v", path, err)
				return err
			}
			if !info.IsDir() {
				// Process YAML file
				if err = processYAMLFile(path, filepath.Join(currentDir, targetFolder, path[len(sourceFolderPath)+1:])); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error traversing folder A: %v", err)
			return err
		}
	}

	fmt.Println("YAML files processed successfully.")
	return nil
}

func getListOfLanguage(dirPath string) (langs []string, err error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Iterate over the fileInfos
	for _, fileInfo := range fileInfos {
		// Check if it's a directory
		if fileInfo.IsDir() {
			langs = append(langs, fileInfo.Name())
		}
	}
	return
}

// processYAMLFile processes a YAML file, merging it with the corresponding file in folder B
func processYAMLFile(pathA, pathB string) error {
	// Read YAML data from file A
	dataA, err := os.ReadFile(pathA)
	if err != nil {
		log.Printf("Failed to read file A %s: %v", pathA, err)
		return err
	}

	// Read YAML data from file B
	dataB, err := os.ReadFile(pathB)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Failed to read file B %s: %v", pathB, err)
			return err
		} else {
			if file, err := os.Create(pathB); err != nil {
				return err
			} else {
				file.Close()
				dataB = nil
			}
		}
	}

	// Unmarshal YAML data from file A
	var yamlDataA interface{}
	err = yaml.Unmarshal(dataA, &yamlDataA)
	if err != nil {
		log.Printf("Failed to unmarshal YAML from file A %s: %v", pathA, err)
		return err
	}

	// Unmarshal YAML data from file B
	var yamlDataB interface{}
	if len(dataB) > 0 {
		err = yaml.Unmarshal(dataB, &yamlDataB)
		if err != nil {
			log.Printf("Failed to unmarshal YAML from file B %s: %v", pathB, err)
			return err
		}
	} else {
		yamlDataB = make(map[string]interface{})
	}

	// Merge YAML data
	mergedYAML := mergeYAML(yamlDataA, yamlDataB)

	// Marshal merged YAML data
	outputData, err := yaml.Marshal(&mergedYAML)
	if err != nil {
		log.Printf("Failed to marshal merged YAML for file %s: %v", pathB, err)
		return err
	}

	// Write to file B
	err = os.WriteFile(pathB, outputData, 0644)
	if err != nil {
		log.Printf("Failed to write to file B %s: %v", pathB, err)
		return err
	}

	fmt.Printf("YAML file %s processed and saved to %s\n", pathA, pathB)
	return nil
}

// mergeYAML merges YAML data from file A into file B
func mergeYAML(dataA, dataB interface{}) interface{} {
	mapA := dataA.(map[string]interface{})
	mapB := dataB.(map[string]interface{})

	for key, value := range mapA {
		// Check if key exists in B
		if _, ok := mapB[key]; !ok {
			mapB[key] = value
		} else {
			// Recursively merge if the value is a map
			if reflect.TypeOf(value).Kind() == reflect.Map {
				mergedValue := mergeYAML(value, mapB[key])
				mapB[key] = mergedValue
			}
		}
	}
	return mapB
}
