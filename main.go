package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type RefResponse struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		Sha  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}

func getRemoteRepoSha() (sha string) {
	headUrl := "https://api.github.com/repos/xionter/go-test-ci-tool/git/ref/heads/main"
	resp, err := http.Get(headUrl)
	check(err)

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	check(err)

	var gitRefResponse RefResponse
	err = json.Unmarshal(body, &gitRefResponse)
	check(err)
	remoteSha := gitRefResponse.Object.Sha
	return remoteSha
}

func main() {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	check(err)
	localSha := strings.TrimSpace(string(out))
	fmt.Printf("сейчас у ветки такой хэш - %v \n", localSha)

	ticker := time.NewTicker(30 * time.Second)
	for range ticker.C {
		remoteSha := getRemoteRepoSha()
		fmt.Printf("сейчас в репозитории на гитхабе такой хэш - %v \n", remoteSha)

		if localSha != remoteSha {
			fmt.Println("новый комит блин")
			break
		}
	}
	ticker.Stop()
}
