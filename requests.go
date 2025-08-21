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


func GetPageContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetOnionContent(url string) ([]byte, error) {
    socksProxy := "127.0.0.1:9050"

    dialer, err := proxy.SOCKS5("tcp", socksProxy, nil, proxy.Direct)
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

    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("bad status: %s", resp.Status)
    }

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetParsedPageContent(url string) (*html.Node, error){
	resp, err := http.Get(url)
	if (err != nil) {
		return &html.Node{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &html.Node{}, err
	}

	return html.Parse(resp.Body)
}

func DownloadDocumentReturnHash(url string, filePath string) (string, error) {
	data, err := GetPageContent(url)

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

	_, err = file.Write([]byte(data))
	if err != nil {
		return "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Convert byte to string
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
