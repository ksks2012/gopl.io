// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Practice 4.14: Create a web server that queries GitHub and generates bug reports, milestones, and corresponding user information
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	github "gopl.io/ch4/githubp"
)

var issueList = template.Must(template.New("issuelist").Funcs(template.FuncMap{
	"repoName": func(repoURL string) string {
		parts := strings.Split(repoURL, "/")
		if len(parts) >= 2 {
			return parts[len(parts)-1]
		}
		return ""
	},
}).Parse(`
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
    <th>#</th>
    <th>State</th>
    <th>User</th>
    <th>Title</th>
    <th>Milestones</th>
</tr>
{{range .Items}}
<tr>
    <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
    <td>{{.State}}</td>
    <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
    <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
    <td>
        <form method="POST" action="/milestones" style="margin:0;">
            <input type="hidden" name="owner" value="{{.User.Login}}">
            <input type="hidden" name="repo" value="{{.RepositoryURL | repoName}}">
            <button type="submit">Show Milestones</button>
        </form>
    </td>
</tr>
{{end}}
</table>
`))

var milestoneList = template.Must(template.New("milestonelist").Funcs(template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("2006-01-02")
	},
}).Parse(`
<h1>Milestones for {{.Owner}}/{{.Repo}}</h1>
<ul>
{{range .Milestones}}
	<li>
		<a href="{{.HTMLURL}}">{{.Title}}</a> (due: {{.DueOn | formatDate}})
	</li>
{{end}}
</ul>
`))

func listIssues(w http.ResponseWriter, r *http.Request) {
	terms := r.URL.Query()["q"]
	if len(terms) == 0 {
		http.Error(w, "Missing search terms. Please provide a 'q' parameter in the URL, e.g., /?q=go+json", http.StatusBadRequest)
		return
	}

	result, err := github.SearchIssues(terms)
	if err != nil {
		log.Printf("Error searching GitHub issues: %v", err)
		http.Error(w, fmt.Sprintf("Error searching issues: %v", err), http.StatusInternalServerError)
		return
	}

	if err := issueList.Execute(w, result); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, fmt.Sprintf("Error rendering results: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Handled request for %s with terms %v. Found %d issues.", r.URL.Path, terms, result.TotalCount)
}

func listMilestones(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	owner := r.FormValue("owner")
	repo := r.FormValue("repo")
	if owner == "" || repo == "" {
		http.Error(w, "Missing owner or repo parameters", http.StatusBadRequest)
		return
	}

	milestones, err := github.SearchMilestones(owner, repo)
	if err != nil {
		log.Printf("Error fetching milestones: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching milestones: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Owner      string
		Repo       string
		Milestones []*github.Milestone
	}{
		Owner:      owner,
		Repo:       repo,
		Milestones: milestones,
	}
	if err := milestoneList.Execute(w, data); err != nil {
		log.Printf("Error executing milestone template: %v", err)
		http.Error(w, fmt.Sprintf("Error rendering milestones: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Handled request for milestones of %s/%s. Found %d milestones.", owner, repo, len(milestones))
}

func main() {
	http.HandleFunc("/", listIssues)
	http.HandleFunc("/milestones", listMilestones)

	port := ":8000"
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
