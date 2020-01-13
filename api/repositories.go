package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"github.com/opendevstack/mockbucket/utils"
	"log"
	"net/http"
	"path"
	"strings"
)

func CreateRepository(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	projectKey := strings.ToLower(vars["projectKey"])
	log.Printf("Project key: %s", projectKey)

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(w)
	var repository Repository
	err := decoder.Decode(&repository)
	if err != nil {
		panic(err)
	}
	repository.Slug = slug.Make(repository.Name)
	log.Println(repository.Slug)

	err = createGitRepository(repository, projectKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(repository)

}

func createGitRepository(repository Repository, project string) error {
	dir := path.Join("/scm", project, fmt.Sprintf("%s.git", repository.Slug))
	stdout, stderr, err := utils.RunCommand("mkdir", "-p", dir)
	if err != nil {
		log.Printf("StdOut:\n%s\nStdErr:\n%s", stdout, stderr)
		return err
	}

	stdout, stderr, err = utils.RunCommand("git", "-C", dir, "init", "--bare")
	if err != nil {
		log.Printf("StdOut:\n%s\nStdErr:\n%s", stdout, stderr)
		return err
	}

	return nil
}
