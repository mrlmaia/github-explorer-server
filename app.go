package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

var baseUrl = "https://api.github.com/"

type Repository struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HtmlUrl     string `json:"html_url"`
}

func repositoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		repositoryHandlerGet(w, r)
	}

}

func repositoryHandlerGet(w http.ResponseWriter, r *http.Request) {
	retryClient := retryablehttp.NewClient()

	name := r.URL.Query().Get("full_name")

	apiResponse, err := retryClient.Get(baseUrl + "repos/" + name)
	log.Println("***********")
	log.Println("StatusCode: ", apiResponse.StatusCode)
	log.Println("ContentLength: ", apiResponse.ContentLength)
	log.Println("***********")

	if err != nil {
		log.Fatalln("Error at the request")
		json.NewEncoder(w).Encode(err)
	}

	defer apiResponse.Body.Close()

	data, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Repository{}

	json.Unmarshal(data, &result)

	json.NewEncoder(w).Encode(result)
}

func makeHttpRequest(url string, data interface{}) {

}

func main() {
	http.HandleFunc("/repository", repositoryHandler)
	http.ListenAndServe(":8080", nil)
}
