package test

import (
	"os"
	"context"
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/Bl4omArchie/simple"
)

// MOCK	for GetContent
func TestGetContentMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hello world"))	}))
	defer ts.Close()

	var ctx context.Context = context.Background()

	content, err := simple.GetContent(ctx, ts.URL, simple.HttpClient(), nil)
	if err != nil {
		t.Fatalf("GetContent failed: %v", err)
	}
	if string(content) != "hello world" {
		t.Errorf("Unexpected content: %s", string(content))
	}
}

// MOCK	for GetContent with custom request
func TestGetContentCustomReqMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hello world"))	}))
	defer ts.Close()

	var ctx context.Context = context.Background()

	content, err := simple.GetContent(ctx, ts.URL, simple.HttpClient(), func(req *http.Request) error {
		req.Header.Set("User-Agent", "Go-http-client/1.1")
		req.Header.Set("Accept", "*/*")
		return nil
	})
	if err != nil {
		t.Fatalf("GetContent with custom req failed: %v", err)
	}
	if string(content) != "hello world" {
		t.Errorf("Unexpected content: %s", string(content))
	}
}

// MOCK for DownloadDocument
func TestDownloadDocumentMock(t *testing.T) {
	var ctx context.Context = context.Background()
	var url string = "https://freetestdata.com/document-files/pdf/"

	var folder string = "doc/"
	var path string = "doc/test.pdf"

	err := simple.DownloadDocument(ctx, url, path, simple.HttpClient(), nil)
	if err != nil {
		t.Fatalf("GetContent failed: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("couldn't find downloaded document")
	}

	err = os.RemoveAll(folder)
	if err != nil {
		t.Fatalf("couldn't delete file")
	}
}
