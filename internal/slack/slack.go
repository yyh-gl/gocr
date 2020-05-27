package slack

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
)

type (
	Client struct {
		webHook   string
		channel   string
		username  string
		iconEmoji string
	}

	PullRequest struct {
		Title       string
		URL         string
		Reviewers   []string
		IsMergeable bool
	}
)

func NewClient(wh, c, un, ie string) *Client {
	return &Client{
		webHook:   wh,
		channel:   c,
		username:  un,
		iconEmoji: ie,
	}
}

func (c Client) Send(t string) {
	// TODO: load from config file
	color := "warning"
	a := slack.Attachment{Color: &color}

	p := slack.Payload{
		Username:    c.username,
		IconUrl:     "",
		IconEmoji:   ":" + c.iconEmoji + ":",
		Channel:     c.channel,
		Text:        t,
		LinkNames:   "true",
		Attachments: []slack.Attachment{a},
		UnfurlLinks: false,
		UnfurlMedia: false,
		Markdown:    false,
	}

	slack.Send(c.webHook, "", p)
}

func CreateMessage(repoName string, prs []PullRequest, userMap []string) string {
	um := make(map[string]string, 0)
	if len(userMap) > 0 {
		for _, m := range userMap {
			splitUserMap := strings.SplitN(m, ":", 2)
			um[splitUserMap[0]] = splitUserMap[1]
		}
	}

	msg := "▼ *" + repoName + "*\n"
	for i, pr := range prs {
		mentions := make([]string, len(pr.Reviewers))
		for i, r := range pr.Reviewers {
			if len(userMap) > 0 {
				mentions[i] = "@" + um[r]
			} else {
				mentions[i] = "@" + r
			}
		}

		tmp := make([]string, 4)
		tmp[0] = fmt.Sprintf("\n*%d: %s*", i+1, pr.Title)
		tmp[1] = pr.URL
		tmp[2] = strings.Join(mentions, ", ")
		tmp[3] = "Please review"
		if pr.IsMergeable {
			tmp[3] = "Let's Merge!"
		}
		msg += strings.Join(tmp, "\n") + "\n"
	}
	return msg
}
