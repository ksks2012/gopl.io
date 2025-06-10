// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package github

import (
	"time"
)

const (
	IssuesURL     = "https://api.github.com/search/issues"
	MilestonesURL = "https://api.github.com/repos/%s/%s/milestones"
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number        int
	HTMLURL       string `json:"html_url"`
	Title         string
	State         string
	User          *User
	CreatedAt     time.Time `json:"created_at"`
	Body          string    // in Markdown format
	RepositoryURL string    `json:"repository_url"` // Added to match GitHub API
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Milestone struct {
	Title   string    `json:"title"`
	DueOn   time.Time `json:"due_on"`
	HTMLURL string    `json:"html_url"`
}
