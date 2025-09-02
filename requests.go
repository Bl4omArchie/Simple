package simple

import (
	"os"
	"io"
	"fmt"
	"time"
	"net/http"
	"crypto/sha256"
	"path/filepath"

	"golang.org/x/net/html"
	"golang.org/x/net/proxy"
)


// Get the body of a webpage.
// Give the url and it will return the html body
func GetContent(url string, client *http.Client) ([]byte, error) {
	body, err := fetchBody(url, client)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return io.ReadAll(body)
}

// Get the body and parse it.
// The package net/html allow to parse a html body into nodes to easily retrieve every html tags
func GetParsedContent(url string, client *http.Client) (*html.Node, error) {
	body, err := fetchBody(url, client)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return html.Parse(body)
}

// Core function to get a webpage body 
func fetchBody(url string, client *http.Client) (io.ReadCloser, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// http client
// Default config
func HttpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

// onion client
// Default config
func OnionClient() (*http.Client, error) {
    dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
    if err != nil {
        return nil, err
    }

    transport := &http.Transport{
        Dial: dialer.Dial,
        DialContext: nil,
    }

    client := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }

	return client, nil
}

// download a document
// give the url, the filepath where you want to store the document and the client
func DownloadDocument(url string, filePath string, client *http.Client) error {
	body, err := fetchBody(url, client)
	if err != nil {
		return fmt.Errorf("failed to fetch body: %w", err)
	}
	defer body.Close()

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// download a document and return its hash value 
func DownloadDocumentReturnHash(url string, filePath string, client *http.Client) (string, error) {
	body, err := fetchBody(url, client)
	if err != nil {
		return "", fmt.Errorf("failed to fetch body: %w", err)
	}
	defer body.Close()

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directories: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)

	if _, err := io.Copy(writer, body); err != nil {
		return "", fmt.Errorf("failed to write file and compute hash: %w", err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
