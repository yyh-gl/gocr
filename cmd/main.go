package main

import (
	"fmt"

	"github.com/yyh-gl/gocr/internal/github"
	"github.com/yyh-gl/gocr/internal/slack"
	"github.com/yyh-gl/gocr/internal/yaml"
)

func main() {
	ct := yaml.LoadConfigFile()

	for _, r := range ct.Repositories {
		//c := github.NewEnterpriseClient(r.EnterpriseHost, r.AccessToken)
		c := github.NewGeneralClient(r.AccessToken)
		prs, err := c.FetchPullRequestDetails(r.Owner, r.Name)
		if err != nil {
			fmt.Println(err)
		}

		if len(*prs) > 0 {
			ss := ct.Slacks[r.SlackID]
			sc := slack.NewClient(ss.WebHook, ss.Channel, ss.Username, ss.IconEmoji)

			msg := slack.CreateMessage(r.Name, prs.ConvertToSlackDTOs(), ss.UserMap)
			sc.Send(msg)
		}
	}
}
