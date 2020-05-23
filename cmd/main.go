package main

import (
	"fmt"

	"github.com/yyh-gl/pr-review-notification/internal/github"
	"github.com/yyh-gl/pr-review-notification/internal/yaml"
)

func main() {
	ct := yaml.LoadRepositoryConfig()

	for repoName, repoDetail := range ct.Repositories {
		c := github.NewEnterpriseClient(repoDetail.EnterpriseHost, repoDetail.AccessToken)
		prs, err := c.FetchPullRequestDetails(repoDetail.Owner, repoName)
		if err != nil {
			fmt.Println(err)
		}

		for _, pr := range prs {
			fmt.Println(pr)
		}
	}
}
