package main

import (
	"fmt"

	"github.com/yyh-gl/gocr/internal/github"
	"github.com/yyh-gl/gocr/internal/yaml"
)

func main() {
	ct := yaml.LoadRepositoryConfig()

	for _, r := range ct.Repositories {
		c := github.NewEnterpriseClient(r.EnterpriseHost, r.AccessToken)
		prs, err := c.FetchPullRequestDetails(r.Owner, r.Name)
		if err != nil {
			fmt.Println(err)
		}

		for _, pr := range prs {
			fmt.Println(pr)
		}
	}
}
