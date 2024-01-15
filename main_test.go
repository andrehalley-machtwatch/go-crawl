package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestReadUrls(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test-urls.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write sample URLs to the temporary file
	urls := []string{
		"https://google.com",
		"https://go.dev",
	}
	for _, u := range urls {
		_, err := tmpFile.WriteString(u + "\n")
		if err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
	}

	// Read URLs from the temporary file using the function being tested
	readURLs, err := readFileUrls(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading URLs from file: %v", err)
	}

	// Check if the read URLs match the expected URLs
	for i, u := range urls {
		if readURLs[i] != u {
			t.Errorf("Expected URL: %s, Got: %s", u, readURLs[i])
		}
	}
}

func TestSaveCrawlResult(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mock HTML Content"))
	}))
	defer mockServer.Close()

	// Call the function being tested
	err := saveCrawlResult(mockServer.URL)
	if err != nil {
		t.Fatalf("Error crawling and saving: %v", err)
	}

	// Check if the file is created and contains the expected content
	parseUrl, err := url.Parse(mockServer.URL)
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}
	filePath := filepath.Join("result", parseUrl.Hostname()+".html")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	expectedContent := "Mock HTML Content"
	if string(content) != expectedContent {
		t.Errorf("Expected content: %s, Got: %s", expectedContent, string(content))
	}
}

// Add more tests as needed
