package slack

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

type (
	Client struct {
		webHook  string
		username string
	}

	Payload struct {
		fromRepoName  string
		toChannelName string
		text          string
	}

	PullRequest struct {
		Title       string
		URL         string
		Reviewers   []string
		IsMergeable bool
	}
)

func NewClient(wh, un string) *Client {
	return &Client{
		webHook:  wh,
		username: un,
	}
}

func (c Client) Send(p Payload) {
	sp := slack.Payload{
		Username:    c.username,
		IconUrl:     "",
		IconEmoji:   "",
		Channel:     "",
		Text:        p.text,
		LinkNames:   "",
		Attachments: nil,
		UnfurlLinks: false,
		UnfurlMedia: false,
		Markdown:    false,
	}

	slack.Send(c.webHook, "", sp)
}

func CreateMessage(prs []PullRequest) string {
	msg := ""
	for i, pr := range prs {
		tmp := make([]string, 4)
		tmp[0] = fmt.Sprintf("*%d: %s*", i+1, pr.Title)
		tmp[1] = pr.URL
		tmp[2] = "@" + strings.Join(pr.Reviewers, ", @")
		tmp[3] = "Please review"
		if pr.IsMergeable {
			tmp[3] = "Let's Merge!"
		}
		msg += strings.Join(tmp, "\n")
	}
	return msg
}
