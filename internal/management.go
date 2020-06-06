package internal

import (
	"context"
	"net/http"

	"github.com/yyh-gl/gocr/internal/repository"
	"github.com/yyh-gl/gocr/internal/repository/github"
	"github.com/yyh-gl/gocr/internal/sender"
	"github.com/yyh-gl/gocr/internal/sender/slack"
)

const (
	RepositoryVendorGitHub = "github"

	SenderVendorSlack = "slack"
)

type Manager struct {
	repositories map[string]repository.Repository
	senders      map[string]sender.Sender
}

func NewManger(configPath string) *Manager {
	config := loadConfigFile(configPath)

	hc := http.DefaultClient

	repos := make(map[string]repository.Repository)
	for vendor, rs := range config.RepositorySet {
		switch vendor {
		case RepositoryVendorGitHub:
			for id, r := range rs {
				repos[id] = github.NewClient(hc, r)
			}
		}
	}

	senders := make(map[string]sender.Sender)
	for vendor, ss := range config.SenderSet {
		switch vendor {
		case SenderVendorSlack:
			for id, s := range ss {
				senders[id] = slack.NewClient(s)
			}
		}
	}

	return &Manager{
		repositories: repos,
		senders:      senders,
	}
}

func (m Manager) Do(ctx context.Context) error {
	for _, r := range m.repositories {
		crs, err := r.FetchCodeReviewRequests(ctx)
		if err != nil {
			return err
		}

		msgs, err := crs.ConvertToMsgMaterials()
		if err != nil {
			return err
		}

		if err := m.senders[r.SenderID()].Send(r.Name(), msgs); err != nil {
			return err
		}
	}

	return nil
}
