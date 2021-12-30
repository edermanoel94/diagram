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

	log.Printf("Create temporary file in %s\n", tempFile.Name())

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

	var sequenceDiagramaResponse SequenceDiagramResponse

	if err := json.Unmarshal(bytesResponse, &sequenceDiagramaResponse); err != nil {
		return nil, err
	}

	log.Printf("Response with image partial url: %s \n", sequenceDiagramaResponse.Img)

	return &sequenceDiagramaResponse, nil
}
