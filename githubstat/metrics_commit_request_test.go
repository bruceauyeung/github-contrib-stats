package githubstat

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"testing"
)

// var proxyClient = &ProxyClient{}
// var client = proxyClient.getClient()
var owner = "kubernetes"
var repo = "kubernetes.github.io"
var author = "xialonglee"

//e.g. https://api.github.com/search/issues?page=1&per_page=100&q=[SHA]+repo:kubernetes/kubernetes.github.io+type:pr+author:[author]
func Test_getStackalyticsCommits(t *testing.T) {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: "put your token to test"},
	}
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	wrapRepositoryCommits := getStackalyticsCommits(client, owner, repo, author)
	for _, c := range wrapRepositoryCommits {
		fmt.Printf("commit message is: %s \n", *c.RepositoryCommit.Commit.Message)
		fmt.Printf("merge time is: %s \n", c.MergedAt.String())
	}
}
