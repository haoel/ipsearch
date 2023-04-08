package ipsearch_test

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/haoel/ipsearch"
	"github.com/stretchr/testify/assert"
)

var (
	FileLines = getFilesLineNum(IPv4CIDRFile)
)

// Get line number of file
func getFilesLineNum(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	var lines int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}
	return lines
}

func TestFile(t *testing.T) {
	lines, err := ipsearch.ReadFile("not_exist_file")
	assert.NotNil(t, err)
	assert.Nil(t, lines)

	lines, err = ipsearch.ReadFile(IPv4CIDRFile)
	test(t, lines, err)
}

func TestURL(t *testing.T) {
	// Connect to a non-existent server
	url := "http://127.0.0.1:9999"
	lines, err := ipsearch.ReadFileFromURL(url)
	assert.NotNil(t, err)
	assert.Nil(t, lines)

	// Start a local server
	go func() {
		fs := http.FileServer(http.Dir(TestDir))
		http.ListenAndServe("127.0.0.1:9999", fs)
	}()

	// Connect to a non-existent file
	url = "http://127.0.0.1:9999/not_exist_file"
	for {
		lines, err = ipsearch.ReadFileFromURL(url)
		if strings.Contains(err.Error(), "connection refused") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		assert.NotNil(t, err)
		assert.Nil(t, lines)
		break
	}

	// Connect to a exist file
	url = "http://127.0.0.1:9999/china_ip_list.txt"
	lines, err = ipsearch.ReadFileFromURL(url)
	test(t, lines, err)

}

func test(t *testing.T, cidrs []string, err error) {
	assert.Nil(t, err)
	assert.Equal(t, len(cidrs), FileLines)
}
