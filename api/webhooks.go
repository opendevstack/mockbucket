package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
)

type WebHook struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	CreatedDate   int64    `json:"createdDate"`
	UpdatedDate   int64    `json:"updatedDate"`
	Events        []string `json:"events"`
	Configuration struct {
		Secret string `json:"secret"`
	} `json:"configuration"`
	URL    string `json:"url"`
	Active bool   `json:"active"`
}

type RefChangedPayload struct {
	EventKey string `json:"eventKey"`
	Date     string `json:"date"`
	Actor    struct {
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
		ID           int    `json:"id"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
		Slug         string `json:"slug"`
		Type         string `json:"type"`
	} `json:"actor"`
	Repository struct {
		Slug          string `json:"slug"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		ScmID         string `json:"scmId"`
		State         string `json:"state"`
		StatusMessage string `json:"statusMessage"`
		Forkable      bool   `json:"forkable"`
		Project       struct {
			Key    string `json:"key"`
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Public bool   `json:"public"`
			Type   string `json:"type"`
		} `json:"project"`
		Public bool `json:"public"`
	} `json:"repository"`
	Changes []struct {
		Ref struct {
			ID        string `json:"id"`
			DisplayID string `json:"displayId"`
			Type      string `json:"type"`
		} `json:"ref"`
		RefID    string `json:"refId"`
		FromHash string `json:"fromHash"`
		ToHash   string `json:"toHash"`
		Type     string `json:"type"`
	} `json:"changes"`
}

func CreateWebhook(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	projectKey := strings.ToLower(vars["projectKey"])
	repositorySlug := strings.ToLower(vars["repositorySlug"])
	log.Printf("Project key: %s", projectKey)

	decoder := json.NewDecoder(req.Body)
	encoder := json.NewEncoder(w)
	var webHook WebHook
	err := decoder.Decode(&webHook)
	if err != nil {
		panic(err)
	}
	_ = createWebkook(webHook, projectKey, repositorySlug)
	w.WriteHeader(http.StatusOK)
	_ = encoder.Encode(webHook)
}

func createWebkook(webhook WebHook, projectKey string, repositorySlug string) error {
	dir := path.Join("/scm", projectKey, fmt.Sprintf("%s.git", repositorySlug))

	file := path.Join(dir, "hooks", "post-receive")
	_ = ioutil.WriteFile(file, []byte(fmt.Sprintf(`#!/bin/sh

if [ -n "$1" -a -n "$2" -a -n "$3" ]; then
	# Output to the terminal in command line mode - if someone wanted to
	# resend an email; they could redirect the output to sendmail
	# themselves
	echo "$2 $3 $1"
else
	while read oldrev newrev refname
	do
		prep_for_email $oldrev $newrev $refname || continue
		generate_email $maxlines | send_mail
	done
fi

echo "$@"
# curl -H "Authorization: Basic %s" %s`, webhook.Configuration.Secret, webhook.URL)), 0755)

	return nil
}
