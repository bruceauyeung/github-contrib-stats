package main

import (
	"flag"
	"fmt"
	"strings"

	"time"

	"./githubstat"
)

func main() {
	start := time.Now()
	flagMetrics := flag.String("metrics", "", "available metrics: (pr)")
	dimension := flag.String("dimension", "", "available dimension: (overall)")
	flag.Parse()

	if flagMetrics == nil || *flagMetrics == "" {
		flagMetrics = &githubstat.Config.Metrics
		if flagMetrics == nil || *flagMetrics == "" {
			fmt.Println("metrics not specified.")
		}
	}
	var metricsRequest githubstat.MetricsRequest
	var metricsParameters githubstat.MetricsParameters
	metricsParameters.Dimension = dimension
	switch *flagMetrics {
	case "issue":
		metricsRequest = &githubstat.IssueMetricsRequest{}
	case "pr":
		metricsRequest = &githubstat.PullRequestMetricsRequest{}
	default:
		metricsRequest = &githubstat.DefaultMetricsRequest{}
	}

	parameters := flag.Args()

	if len(parameters) == 0 {
		parameters = githubstat.Config.Repos

	}
	for _, repoStr := range parameters {
		repo := strings.Split(repoStr, "/")
		if len(repo) != 2 {
			fmt.Printf("invalid repository name : %s, must be of format 'ownername/reponame'\n", repoStr)
		}
		metricsParameters.Repos = append(metricsParameters.Repos,
			&githubstat.RepoParameters{&repo[0], &repo[1]})
	}

	metricsRequest.SetParameters(&metricsParameters)
	metrics := metricsRequest.FetchMetrics()
	metrics.Show()
	elapsed := time.Since(start)
	fmt.Printf("stats finished and spent %v minutes", elapsed.Minutes())
}
