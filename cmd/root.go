package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yyh-gl/gocr/internal"
)

var (
	configPath string

	rootCmd = &cobra.Command{
		Use:     "gocr",
		Version: "0.3.3",
		Short:   "GoCR is planner for a notification of code review request.",
		Long:    "GoCR provides easy way for notifying request of code review.\nYou will soon be able to start a notification of code review request.",
		Run: func(cmd *cobra.Command, args []string) {
			m := internal.NewManger(configPath)
			if err := m.Do(context.Background()); err != nil {
				fmt.Println(err)
				os.Exit(1)
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
