package main

import (
	"context"
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

	for _, issue := range issues {
		found := false
		for _, i := range myIssues {
			if *i.Title == *issue.Title {
				found = true
				return
			}
		}
		if !found {
			log.Infof("Creating new issue %v", *issue.Title)
			_, _, err := c.Issues.Create(context.Background(), cfg.RepoOwner, cfg.RepoName, &issue)
			if err != nil {
				CheckErr(err, "Fail to create issue", -1)
			}
		}
	}
}

func getIssues(client *github.Client, cfg *Config) (issues []*github.Issue, err error) {
	page := 1
	for page != 0 {
		lsopt := &github.ListOptions{
			Page: page,
		}
		opt := &github.IssueListByRepoOptions{
			ListOptions: *lsopt,
			State:       "all",
		}
		is, resp, err := client.Issues.ListByRepo(context.Background(), cfg.RepoOwner, cfg.RepoName, opt)
		if err != nil {
			log.Fatal(err)
		}
		page = resp.NextPage
		log.Infof("Adding %v issues", len(is))
		issues = append(issues, is...)
	}

	return
}
