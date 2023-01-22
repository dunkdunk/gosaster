package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

func visit(path string, di fs.DirEntry, err error) error {
	// If the path is a file, exec grep on it
	if !di.IsDir() {
		// Grep the file for the test string "<p>" and return the results with a line number and colorized
		cmd := exec.Command("grep", "-n", "--color=always", "<p*>", path)

		// Get the output of the command
		out, err := cmd.Output()

		// If the output is code 0, print the output
		if err == nil {
			fmt.Printf("%s", out)
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

	// Print the root directory
	fmt.Printf("Starting scan in %s\n", root)

	err := filepath.WalkDir(root, visit)

	if err != nil {
		fmt.Println(err)
	}
}
