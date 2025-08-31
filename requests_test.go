package simple

import (
	"testing"
	"net/http"
	"net/http/httptest"
)


// MOCK for default
func TestGetContentMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hello world"))	}))
	defer ts.Close()

	client := HttpClient()
	content, err := GetContent(ts.URL, client)
	if err != nil {
		t.Fatalf("GetContent failed: %v", err)
	}
	if string(content) != "hello world" {
		t.Errorf("Unexpected content: %s", string(content))
	}
}
