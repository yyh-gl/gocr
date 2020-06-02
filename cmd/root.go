package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/yyh-gl/gocr/internal/github"
	"github.com/yyh-gl/gocr/internal/slack"
	"github.com/yyh-gl/gocr/internal/yaml"
)

var (
	configPath string

	rootCmd = &cobra.Command{
		Use:   "gocr",
		Short: "GoCR is simple code review notifier",
		// TODO
		Long: `WIP`,
		Run: func(cmd *cobra.Command, args []string) {
			ct := yaml.LoadConfigFile(configPath)

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
		},
	}
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.Flags().StringVarP(&configPath, "cfgPath", "c", homeDir+"/.gocr.yml", "Path to config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
