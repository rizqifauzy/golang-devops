package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v45/github"
	"k8s.io/client-go/kubernetes"
)

type server struct {
	client           *kubernetes.Clientset
	githubClient     *github.Client
	webhookSecretKey string
}

func (s server) webhook(w http.ResponseWriter, req *http.Request) {
	payload, err := github.ValidatePayload(req, []byte(s.webhookSecretKey))
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("ValidatePayload Error: %s\n", err)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		w.WriteHeader(500)
		fmt.Printf("ParseWebHook Error: %s\n", err)
		return
	}
	switch event := event.(type) {
	case *github.PushEvent:
		files := getFiles(event.Commits)
		fmt.Printf("Found files: %s\n", strings.Join(files, ", "))
	default:
		w.WriteHeader(500)
		fmt.Printf("Event not found: %s\n", event)
		return
	}
}

func getFiles(commits []*github.HeadCommit) []string {
	allFiles := []string{}
	for _, commit := range commits {
		allFiles = append(allFiles, commit.Added...)
		allFiles = append(allFiles, commit.Modified...)
	}
	allUniqueFiles := make(map[string]bool)
	for _, filename := range allFiles {
		allUniqueFiles[filename] = true
	}
	allUniqueFilesSlice := []string{}
	for filename := range allUniqueFiles {
		allUniqueFilesSlice = append(allUniqueFilesSlice, filename)
	}
	return allUniqueFilesSlice
}
