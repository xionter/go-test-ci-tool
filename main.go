package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	githubAPIURL   = "https://api.github.com/repos/xionter/go-test-ci-tool/git/ref/heads/"
	workingBranch  = "main"
	testingRepoDir = "test"
	repoName = "go-test-ci-tool"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type RefResponse struct {
	Object struct {
		Sha string `json:"sha"`
	} `json:"object"`
}

func getRemoteSHA() (sha string) {
	headUrl := fmt.Sprintf("%s%s", githubAPIURL, workingBranch)
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

func git(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	out, err := cmd.CombinedOutput()
	check(err)
	return strings.TrimSpace(string(out)), nil
}

func syncTestRepo() (string, error) {
	homeDir, err := os.UserHomeDir()
	check(err)
	testingRepoPath := filepath.Join(homeDir, testingRepoDir, repoName)

	localSHA, err := git("", "rev-parse", "HEAD")
	check(err)
	remoteSHA := getRemoteSHA()
	fmt.Println(localSHA, remoteSHA)
	if _, err := git(testingRepoPath, "fetch", "origin", workingBranch); err != nil {
		return "", err
	}

	if localSHA == remoteSHA {
		return "", nil
	}

	if _, err := git(testingRepoPath, "checkout", workingBranch); err != nil {
		return "", err
	}
	if _, err := git(testingRepoPath, "reset", "--hard", remoteSHA); err != nil {
		return "", err
	}
	return remoteSHA, nil
}

func main() {
	ticker := time.NewTicker(30 * time.Second)

	sha, err := syncTestRepo()
	check(err)
	if sha != "" {
		fmt.Println("новые комиты, запускаю для тестов...")
		//запустить тесты ~test/repo
	} else {
		fmt.Println("пока что не было новых комитов")
	}

	for range ticker.C {
		sha, err := syncTestRepo()
		check(err)
		if sha != "" {
			fmt.Println("новые комиты, запускаю для тестов...")
			//запустить тесты ~test/repo
		} else {
			fmt.Println("пока что не было новых комитов")
		}
	}
	ticker.Stop()
}
