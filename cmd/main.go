package main

import (
	"fmt"
	"net/http"

	"github.com/yyh-gl/gocr/internal/github"
	"github.com/yyh-gl/gocr/internal/slack"
	"github.com/yyh-gl/gocr/internal/yaml"
)

func main() {
	ct := yaml.LoadConfigFile()

	for _, r := range ct.Repositories {
		var c *github.Client
		if r.IsEnterprise {
			c = github.NewEnterpriseClient(http.DefaultClient, r.EnterpriseHost, r.AccessToken)
		} else {
			c = github.NewGeneralClient(http.DefaultClient, r.AccessToken)
		}

		prs, err := c.FetchPullRequestDetails(r.Owner, r.Name)
		if err != nil {
			fmt.Println(err)
		}

		if len(*prs) > 0 {
			ss := ct.Slacks[r.SlackID]
			sc := slack.NewClient(ss.WebHook, ss.Channel, ss.Username, ss.IconEmoji)

			msg := slack.CreateMessage(prs.ConvertToSlackDTOs(), ss.UserMap)
			sc.Send(r.Name, msg)
		}
	}
}
