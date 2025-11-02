package simple

import (
	"os"
	"io"
	"net"
	"fmt"
	"time"
	"strings"
	"context"
	"net/http"
	"crypto/sha256"
	"path/filepath"

	"golang.org/x/net/html"
	"golang.org/x/net/proxy"
)

type CustomReq func(*http.Request) error

// Get the body of a webpage.
// Give the url and it will return the html body
func GetContent(ctx context.Context, url string, client *http.Client, req CustomReq) ([]byte, error) {
	body, err := fetchBody(ctx, url, "get", client, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return io.ReadAll(body)
}

// Get the body and parse it.
// The package net/html allow to parse a html body into nodes to easily retrieve every html tags
func GetParsedContent(ctx context.Context, url string, client *http.Client, req CustomReq) (*html.Node, error) {
	body, err := fetchBody(ctx, url, "get", client, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return html.Parse(body)
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
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
    }

    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }, nil
}

// download a document
// give the url, the filepath where you want to store the document and the client
func DownloadDocument(ctx context.Context, url string, filePath string, client *http.Client, req CustomReq) error {
	body, err := fetchBody(ctx, url, "get", client, req)
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
func DownloadDocumentReturnHash(ctx context.Context, url string, filePath string, client *http.Client, req CustomReq) (string, error) {
	body, err := fetchBody(ctx, url, "get", client, req)
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

// Core function to get a webpage body
// In case of error 400, see this page : https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/User-Agent
// The error can occur if the requested server blocks some headers, use the customReq in order to solve that.
func fetchBody(ctx context.Context, url string, requestType string, client *http.Client, customReq func(*http.Request) error) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(requestType), url, nil)
    if err != nil {
        return nil, err
    }

	if customReq == nil {
		customReq = func(r *http.Request) error {
			getDefaultCustomReq(r)
			return nil
		}
	}
	if err := customReq(req); err != nil {
		return nil, fmt.Errorf("incorrect custom request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// getDefaultCustomReq applies safe headers to reduce 400 or 403 errors.
func getDefaultCustomReq(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5993.90 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
}
