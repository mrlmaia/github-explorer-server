package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	if err != nil {
		w.WriteHeader(500)
		log.Fatalln("Error at the request")
		erro := AppError{Message: "Internal server error"}
		json.NewEncoder(w).Encode(erro)
		return
	}

	if code := apiResponse.StatusCode; code > 299 || code < 200 {
		erro := AppError{}
		if code == 404 {
			erro = AppError{Message: "Repository not found"}
		} else {
			erro = AppError{Message: "An error with GitHub's API has ocurred"}
		}

		w.WriteHeader(400)
		json.NewEncoder(w).Encode(erro)
		return
	}

	defer apiResponse.Body.Close()

	data, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		erro := AppError{Message: "Internal server error"}
		w.WriteHeader(500)
		log.Fatal("Error processing result")
		json.NewEncoder(w).Encode(erro)
	}

	result := Repository{}

	json.Unmarshal(data, &result)

	json.NewEncoder(w).Encode(result)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/repository", repositoryHandlerGet)
	http.ListenAndServe(":8080", router)
}
