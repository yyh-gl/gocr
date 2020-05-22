package main

import (
	"fmt"

	"github.com/yyh-gl/pr-review-notification/prn"
)

func main() {
	//c := prn.NewGeneralClient("")
	c := prn.NewEnterpriseClient()
	prs, err := c.FetchPRs("dmm-app", "pointclub-android")
	if err != nil {
		fmt.Println(err)
	}

	for _, pr := range prs {
		fmt.Println("========================")
		fmt.Println(pr)
		fmt.Println("========================")
	}
}
