package internal

import (
	"context"
	"errors"
	"github.com/google/go-github/v28/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	PR_COMMENT_TITLE = "### Qiniu Code Static Anlysis Result"
)

type GitHubClient struct {
	Client     *github.Client
	PullNumber string
	RepoOwner  string
	RepoName   string
	JobName    string
	PR         *github.PullRequest
}

func NewGitHubClient(path string) (*GitHubClient, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Fail to open the github oauth token file, the error is %v", err)
		return nil, err
	}
	token := string(b)
	token = strings.TrimSpace(token)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tc := oauth2.NewClient(ctx, ts)
	client := GitHubClient{Client: github.NewClient(tc)}
	err = client.getEnvironmentVariables()
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *GitHubClient) getEnvironmentVariables() error {
	c.PullNumber = os.Getenv("PULL_NUMBER")
	c.RepoOwner = os.Getenv("REPO_OWNER")
	c.RepoName = os.Getenv("REPO_NAME")
	c.JobName = os.Getenv("JOB_NAME")

	if c.PullNumber == "" ||
		c.RepoName == "" ||
		c.RepoOwner == "" ||
		c.JobName == "" {
		log.Errorf("Fail to get environment variables, check the following variables to see if they are right. "+
			"PULL_NUMBER: %v, REPO_OWNER: %v, REPO_NAME： %v, JOB_NAME: %v", c.PullNumber, c.RepoOwner, c.RepoName, c.JobName)
		return errors.New("Some environment variables missing.")
	}
	return nil
}

func (c *GitHubClient) GetMyRepoPullRequest() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pullNum, _ := strconv.Atoi(c.PullNumber)
	pullRequest, _, err := c.Client.PullRequests.Get(ctx, c.RepoOwner, c.RepoName, pullNum)
	if err != nil {
		log.Errorf("Fail to get repository from GitHub, the error is %v", err)
		return err
	}
	c.PR = pullRequest
	return nil
}

func (c *GitHubClient) DeletePreviousComments() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opt := &github.PullRequestListCommentsOptions{}
	pullNum, _ := strconv.Atoi(c.PullNumber)
	comments, _, err := c.Client.PullRequests.ListComments(ctx, c.RepoOwner, c.RepoName, pullNum, opt)
	if err != nil {
		log.Errorf("Fail to get repository from GitHub, the error is %v", err)
		return err
	}
	var deletedComments []int64
	for _, comment := range comments {
		if strings.Contains(comment.GetBody(), PR_COMMENT_TITLE) {
			deletedComments = append(deletedComments, comment.GetID())
		}
	}

	for _, comment := range deletedComments {
		_, err := c.Client.PullRequests.DeleteComment(ctx, c.RepoOwner, c.RepoName, comment)
		log.Errorf("Fail to delete the comment, the error is %v", err)
	}

	return nil
}

func (c *GitHubClient) PostCommentsToPR() {

}
