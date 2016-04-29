package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/caarlos0/env"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Config struct {
	ApiBaseUrl string `env:"API_BASE_URL"`
	RepoName   string `env:"REPO_NAME"`
	RepoOwner  string `env:"REPO_OWNER"`
	Token      string `env:"TOKEN"`
}

// CheckErr checks for error, logs and optionally exits the program
func CheckErr(err error, msg string, exit int) {
	if err != nil {
		log.Errorf(msg, err)

		if exit != -1 {
			os.Exit(exit)
		}
	}
}

func main() {
	cfg := &Config{}
	env.Parse(cfg)

	var issues []github.IssueRequest
	bytes, _ := ioutil.ReadAll(os.Stdin)
	json.Unmarshal(bytes, &issues)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	c := github.NewClient(tc)
	c.BaseURL, _ = url.Parse(cfg.ApiBaseUrl)

	myIssues, err := getIssues(c, cfg)
	if err != nil {
		CheckErr(err, "Fail to fetch list of issues", 1)
	}
	issuesMap := make(map[string]github.Issue)
	for _, i := range myIssues {
		issuesMap[*i.Title] = i
	}

	for _, issue := range issues {
		if _, ok := issuesMap[*issue.Title]; !ok {
			log.Info("Creating new issue " + *issue.Title)
			_, _, err := c.Issues.Create(cfg.RepoOwner, cfg.RepoName, &issue)
			if err != nil {
				CheckErr(err, "Fail to create issue", -1)
			}
		}
	}
}

func getIssues(client *github.Client, cfg *Config) (issues []github.Issue, err error) {
	/*page := 1
	for page != 0 {
		lsopt := &github.ListOptions{
			PerPage: 100,
			Page:    page,
		}
		opt := &github.IssueListByRepoOptions{
			ListOptions: *lsopt,
		}
		is, resp, _ := client.Issues.ListByRepo(cfg.RepoOwner, cfg.RepoName, opt)
		page = resp.NextPage
		issues = append(issues, is...)
	}*/

	issues, _, err = client.Issues.ListByRepo(cfg.RepoOwner, cfg.RepoName, &github.IssueListByRepoOptions{})
	log.Printf("Found %v issues", len(issues))

	return
}
