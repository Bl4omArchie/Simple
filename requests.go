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


func GetContent(url string, client *http.Client) ([]byte, error) {
	body, err := FetchBody(url, client)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return io.ReadAll(body)
}

func GetParsedContent(url string, client *http.Client) (*html.Node, error) {
	body, err := FetchBody(url, client)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return html.Parse(body)
}

func FetchBody(url string, client *http.Client) (io.ReadCloser, error) {
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

func HttpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func OnionClient(socksProxy *string) (*http.Client, error) {
	if socksProxy == nil {
        defaultProxy := "127.0.0.1:9050"
        socksProxy = &defaultProxy
	}

    dialer, err := proxy.SOCKS5("tcp", *socksProxy, nil, proxy.Direct)
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

func DownloadDocumentReturnHash(url string, filePath string, client *http.Client) (string, error) {
	body, err := FetchBody(url, client)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	writer := io.MultiWriter(file, hasher)

	if _, err := io.Copy(writer, body); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

