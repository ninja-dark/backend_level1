package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"homework4/model"
)

type Filer interface{
	GetList(extension string) ([]model.File, error)
} 



type Handler struct{
	File Filer
}

type UploadHandler struct {
	HostAddr string
	UploadDir string 
}


func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	filePath := h.UploadDir + "/" + header.Filename
	
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	
	fileLink := h.HostAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodGet:
		name := r.FormValue("ext")
		file, err := h.File.GetList(name)
		if err != nil{
			http.Error(w, "Unable to read files", http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(file)
		if err != nil {
			http.Error(w, "Unable to encoding", http.StatusInternalServerError)
		}
		default:
			http.Error(w, "Unknown request", http.StatusMethodNotAllowed)
	}
}


func main(){

	handler := &Handler{}
	http.Handle("/", handler)

	srv := &http.Server{
		Addr: ":8000",
		Handler: handler,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srv.ListenAndServe()

	uploadHandler := &UploadHandler{
		HostAddr: "http://localhost:8080",
		UploadDir: "upload",
	}

	http.Handle("/upload", uploadHandler)


	dirToServe := http.Dir(uploadHandler.UploadDir)

	fs := &http.Server{
		Addr: ":8080",
		Handler: http.FileServer(dirToServe),
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
}
	fs.ListenAndServe()

}