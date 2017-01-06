package githubstat

import (
	"fmt"
	"github.com/google/go-github/github"
	"strings"
	"time"
)

type PullRequestCommit struct {
	RepositoryCommit *github.RepositoryCommit
	Owner            string
	Repo             string
	MergedAt         *time.Time
}

// sameCommitter means the pr author was same with the commit author
func (m *PullRequestCommit) findMergedTime(client *github.Client, sameCommitter bool) bool {
	SHA := m.RepositoryCommit.SHA
	owner := m.Owner
	repo := m.Repo
	author := m.RepositoryCommit.Author.Login
	if !sameCommitter {
		author = m.RepositoryCommit.Committer.Login
	}
	// TODO consider using template
	query := *SHA + " repo:" + owner + "/" + repo + " type:pr" + " author:" + *author
	pr := findPRfromCommit(client, query)
	if pr == nil {
		return false
	}
	m.MergedAt = pr.ClosedAt
	return true
}

func listCommits(client *github.Client, owner string, repo string, author string) ([]*github.RepositoryCommit, error) {
	opt := &github.CommitsListOptions{
		Author:      author,
		Until:       Config.StatEndTime,
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
	}

	page := 1
	var allCommits []*github.RepositoryCommit
	for {
		commits, resp, err := client.Repositories.ListCommits(owner, repo, opt)
		if err != nil {
			return nil, err
		}

		allCommits = append(allCommits, commits...)
		fmt.Printf("page:%d fin\n", page)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
		page++
	}

	return allCommits, nil
}

func filterCommits(commits []*github.RepositoryCommit) []*github.RepositoryCommit {
	var filteredCommits []*github.RepositoryCommit
	prefixToFiltered := "Merge branch 'master' into"
	for _, c := range commits {
		if strings.HasPrefix(*c.Commit.Message, prefixToFiltered) {
			continue
		}
		filteredCommits = append(filteredCommits, c)
	}
	return filteredCommits
}

// findPRfromCommit finds a pull request from SHA of commit
// Note: every pull request is an issue, but not every issue is a pull request.
// make sure string "+type:pr" was included in query string
func findPRfromCommit(client *github.Client, query string) *github.Issue {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
			// 100 is so big, it should not return so much results
			// so no need to handle the "NextPage" attr of the response.
		},
	}
	results, _, err := client.Search.Issues(query, opt)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	} else if *results.Total > 1 {
		fmt.Printf("warning: find multiple pull requests from a commit's SHA, we take the first one.\n")
		fmt.Printf("query string is %s \n", query)
		// TODO to sort the issues(prs) according to the pr closed time and the earliest pr should be the first item.
		return &results.Issues[0]
	} else if len(results.Issues) == 0 {
		fmt.Printf("warning: find no pull requests from a commit's SHA, the commit should come from commiter.\n")
		fmt.Printf("query string is %s \n", query)
		return nil
	}
	return &results.Issues[0]
}

func getStackalyticsCommits(client *github.Client, owner string, repo string, author string) []*PullRequestCommit {
	fmt.Printf("%s/%s : listing commits of stackalytics.com style\n", owner, repo)
	overallCommits, err := listCommits(client, owner, repo, author)
	if err != nil {
		panic(err)
	}
	overallCommits = filterCommits(overallCommits)

	var warpOverallCommits []*PullRequestCommit
	for _, commit := range overallCommits {
		warpCommit := &PullRequestCommit{commit, owner, repo, nil}
		if !warpCommit.findMergedTime(client, true) {
			fmt.Printf("could not find pull request which includes this commit:%s \n", *commit.Commit.Message)
			fmt.Printf("change author to the committer:%s \n", *commit.Committer.Login)
			warpCommit.findMergedTime(client, false)
		}
		warpOverallCommits = append(warpOverallCommits, warpCommit)
	}
	return warpOverallCommits
}
