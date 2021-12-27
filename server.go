package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	BaseUrl = "https://www.websequencediagrams.com"
)

const (
	Default    = "default"
	Earth      = "earth"
	Magazine   = "magazine"
	ModernBlue = "modern-blue"
	Mscgen     = "mscgen"
	Napkin     = "napkin"
	Omegapple  = "omegapple"
	Patent     = "patent"
	Qsd        = "qsd"
	Rose       = "rose"
	RoundGreen = "roundgreen"
)

type SequenceDiagramRequest struct {
	Message string `json:"message"`
	Style   string `json:"style"`
	Format  string `json:"format"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

type SequenceDiagramResponse struct {
	Img    string   `json:"img"`
	Errors []string `json:"errors"`
}

func (s SequenceDiagramResponse) ImageUrl() string {
	return fmt.Sprintf("%s/%s", BaseUrl, s.Img)
}

func main() {

	http.HandleFunc("/download", handlerDownload)
	http.HandleFunc("/health", handlerHealthcheck)

	log.Println("Starting server in port: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlerHealthcheck(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"error": "method not allowed"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, `{"live": "ok"}`)
}

func handlerDownload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"error": "method not allowed"}`)
		return
	}

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

	defer func() {
		imgFile.Close()
		os.Remove(imgFile.Name())
	}()

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

func downloadImage(imgUrl, format string) (*os.File, error) {

	res, err := http.Get(imgUrl)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	tempFile, err := ioutil.TempFile("", fmt.Sprintf("sequence_diagram_*.%s", format))

	if err != nil {
		return nil, err
	}

	log.Printf("Create file %s, and download image from %s \n", tempFile.Name(), imgUrl)

	io.Copy(tempFile, res.Body)

	return tempFile, nil
}

func urlValuesFromDiagramSequenceRequest(sequenceDiagramReq SequenceDiagramRequest) url.Values {

	data := url.Values{}

	data.Set("apiVersion", "1")
	data.Set("message", sequenceDiagramReq.Message)
	data.Set("style", sequenceDiagramReq.Style)
	data.Set("format", sequenceDiagramReq.Format)

	return data
}

func callWebSequenceDiagramAPI(sequenceDiagramRequest SequenceDiagramRequest) (*http.Response, error) {

	data := urlValuesFromDiagramSequenceRequest(sequenceDiagramRequest)

	res, err := http.PostForm(fmt.Sprintf("%s/index.php", BaseUrl), data)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func getSequenceDiagram(sequenceDiagramRequest SequenceDiagramRequest) (*SequenceDiagramResponse, error) {

	res, err := callWebSequenceDiagramAPI(sequenceDiagramRequest)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bytesResponse, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	log.Printf("Response from websequencediagrams %s \n", string(bytesResponse))

	var sequenceDiagramaResponse SequenceDiagramResponse

	if err := json.Unmarshal(bytesResponse, &sequenceDiagramaResponse); err != nil {
		return nil, err
	}

	return &sequenceDiagramaResponse, nil
}
