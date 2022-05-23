// Package files implements methods for handling program concerned files.
package files

import (
	"bufio"
	"fmt"
	"os"
)

// ReadLinesFromPath is used to get all data from file of specified path,
// returning a string slice of file data and error.
func ReadLinesFromPath(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	return ReadLinesFromFile(file)
}

// ReadLinesFromFile is used to get all data from specified file,
// returning a string slice of file data and error.
func ReadLinesFromFile(file *os.File) ([]string, error) {
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}

// WriteLinesToFile is used to persist data into specified file.
func WriteLinesToFile(file *os.File, lines *[]string) error {
	w := bufio.NewWriter(file)
	for _, line := range *lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}
