package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Global variables
var definitions []Definition
var filesScanned int = 0

type Definition struct {
	Extensions []string
	Regex      string
}

func importDefinitions() []Definition {
	// Let's first read the `functions.json` file
	content, err := ioutil.ReadFile("./definitions/functions.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Unmarshall the data into `payload`
	var payload []Definition
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return payload
}

func visit(path string, di fs.DirEntry, err error) error {
	// If the path is a file, continue
	if !di.IsDir() {
		// Get the file extension
		extension := filepath.Ext(path)

		// Loop through the definitions and see if the extension matches
		for _, definition := range definitions {
			for _, ext := range definition.Extensions {
				if ext == extension {
					// Increment the files scanned counter
					filesScanned++

					// Grep the individual file using the matching regex and return the results with a line number and colorized
					cmd := exec.Command("grep", "-n", "--color=always", "-E", definition.Regex, path)

					out, err := cmd.Output()

					// If the output is code 0, print the output
					if err == nil {
						// Print the file path and output on two lines
						fmt.Printf("%s\n%s", path, out)
					}
				}
			}
		}
	}

	return nil
}

func main() {
	// Return an error if no arugments are passed
	if len(os.Args) < 2 {
		fmt.Println("Please pass the directory you'd like to scan.")
		os.Exit(1)
	}

	// Get the command line arguments and parse them
	// The first argument is the root directory to scan
	root := os.Args[1]

	// Return an error if the root directory does not exist
	if _, err := os.Stat(root); os.IsNotExist(err) {
		fmt.Printf("The root directory %s does not exist\n", root)
		os.Exit(1)
	}

	// Call importDefinitions() to get the definitions and store them in the global variable
	definitions = importDefinitions()

	// Print the root directory
	fmt.Printf("Starting scan in %s\n", root)

	err := filepath.WalkDir(root, visit)

	// Print the number of files scanned
	fmt.Printf("Scanned %d files\n", filesScanned)

	if err != nil {
		fmt.Println(err)
	}
}
