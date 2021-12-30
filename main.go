package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/download", handlerDownload).Methods(http.MethodPost)
	router.HandleFunc("/health", handlerHealthCheck).Methods(http.MethodGet)

	svr := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(svr.ListenAndServe())
}

func handlerHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, `{"live": "ok"}`)
}

func handlerDownload(w http.ResponseWriter, r *http.Request) {

	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}

	defer r.Body.Close()

	log.Printf("Request body: %s \n", string(requestBody))

	var sequenceDiagramRequest SequenceDiagramRequest

	if err := json.Unmarshal(requestBody, &sequenceDiagramRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}

	sequenceDiagramResponse, err := getSequenceDiagram(sequenceDiagramRequest)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}

	imgFile, err := downloadImage(sequenceDiagramResponse.ImageUrl(), sequenceDiagramRequest.Format)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}

	defer imgFile.Close()

	imgFileInfo, err := imgFile.Stat()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}

	tempFileBuffer := make([]byte, 512)

	imgFile.Read(tempFileBuffer)

	fileContentType := http.DetectContentType(tempFileBuffer)

	contentType := fmt.Sprintf("%s;image.%s", fileContentType, sequenceDiagramRequest.Format)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", imgFileInfo.Size()))

	imgFile.Seek(0, 0)

	io.Copy(w, imgFile)
}
