package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadHandler(t *testing.T) {
	file, _ := os.Open("testfile")
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok!")
	}))
	defer ts.Close()

	uploadHandler := &UploadHandler{
		UploadDir: "upload",
		HostAddr:  ts.URL,
	}
	uploadHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `testfile`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/?name=txt", nil)

	rr := httptest.NewRecorder()
	handler := &Handler{}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `Parsed query-param with key "name": txt`

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
		rr.Body.String(), expected)
		}


}