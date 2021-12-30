package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type issueComment struct {
	ID        int32
	Body      string
	User      string
	CreatedAt time.Time
}

type ghClient struct {
	iClient
}

func (c *ghClient) listIssueComments(org, repo string, number int32) ([]issueComment, error) {
	v, err := c.ListPRComments(org, repo, number)
	if err != nil {
		return nil, err
	}

	r := make([]issueComment, len(v))

	for i := range v {
		item := &v[i]

		ct, _ := time.Parse(time.RFC3339, item.CreatedAt)

		r[i] = issueComment{
			ID:        item.Id,
			Body:      item.Body,
			User:      item.User.GetLogin(),
			CreatedAt: ct,
		}
	}

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].CreatedAt.Before(r[j].CreatedAt)
	})

	return r, nil
}

func (c *ghClient) getCommitHashTree(org, repo, SHA string) (string, error) {
	v, err := c.GetPRCommit(org, repo, SHA)
	if err != nil {
		return "", err
	}

	if v.Commit == nil {
		return "", fmt.Errorf("single commit(%s/%s/%s) data is abnormal: %+v", org, repo, SHA, v)
	}

	return v.Commit.Tree.GetSha(), nil
}

func (gc *ghClient) getChangedFiles(org, repo string, number int32) ([]string, error) {
	changes, err := gc.GetPullRequestChanges(org, repo, number)
	if err != nil {
		return nil, fmt.Errorf("cannot get PR changes for %s/%s#%d", org, repo, number)
	}

	v := make([]string, len(changes))
	for i := range changes {
		v[i] = changes[i].Filename
	}

	return v, nil
}

func normalizeLogin(s string) string {
	return strings.TrimPrefix(strings.ToLower(s), "@")
}
