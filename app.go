package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var baseUrl = "https://api.github.com/"

type Repository struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	HtmlUrl     string `json:"html_url"`
	Url         string `json:"url"`
}

type AppError struct {
	Message string `json:"message"`
}

type Response struct {
	Error AppError   `json:"error"`
	Data  Repository `json:"data"`
}

func repositoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		repositoryHandlerGet(w, r)
	}

}

func repositoryHandlerGet(w http.ResponseWriter, r *http.Request) {
	// retryClient := retryablehttp.NewClient()

	repoOwner := r.URL.Query().Get("repo_owner")
	repoName := r.URL.Query().Get("repo_name")

	if repoName == "" || repoOwner == "" {
		erro := AppError{
			Message: "You must provide a repo_owner and repo_name",
		}
		// erro := errors.New("You must provide a repo_owner and repo_name")
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(erro)
		return
	}

	name := repoOwner + "/" + repoName

	apiResponse, err := http.Get(baseUrl + "repos/" + name)
	log.Println("StatusCode: ", apiResponse.StatusCode)

	if err != nil {
		log.Fatalln("Error at the request")
		// json.NewEncoder(w).Encode(err)
	}
	log.Print("err", err)

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
