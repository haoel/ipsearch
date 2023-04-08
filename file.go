package ipsearch

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ReadFile reads a ip list file and returns a slice of strings, one for each line.
func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readLines(file)
}

// ReadFileFromURL from a URL and returns a slice of strings, one for each line.
func ReadFileFromURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code error: %d", resp.StatusCode)
	}
	return readLines(resp.Body)
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
