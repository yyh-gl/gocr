package main

import (
	"fmt"

	"github.com/yyh-gl/gocr/internal/slack"

	"github.com/yyh-gl/gocr/internal/github"

	"github.com/yyh-gl/gocr/internal/yaml"
)

func main() {
	ct := yaml.LoadRepositoryConfig()

	//sc := slack.NewClient("", "")

	for _, r := range ct.Repositories {
		//c := github.NewEnterpriseClient(r.EnterpriseHost, r.AccessToken)
		c := github.NewGeneralClient(r.AccessToken)
		prs, err := c.FetchPullRequestDetails(r.Owner, r.Name)
		if err != nil {
			fmt.Println(err)
		}

		if len(*prs) > 0 {
			msg := slack.CreateMessage(prs.ConvertToSlackDTOs())
			fmt.Println("========================")
			fmt.Println(msg)
			fmt.Println("========================")
		}
	}
}
