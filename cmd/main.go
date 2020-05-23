package main

import (
	"fmt"

	"github.com/yyh-gl/pr-review-notification/internal/github"
)

func main() {
	//c := prn.NewGeneralClient("")
	c := github.NewEnterpriseClient("", "")
	prs, err := c.FetchPullRequestDetails("dmm-app", "pointclub-android")
	if err != nil {
		fmt.Println(err)
	}

	for _, pr := range prs {
		fmt.Println("========================")
		fmt.Println(pr)
		fmt.Println("========================")
	}
}
