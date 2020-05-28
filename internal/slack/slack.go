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
		HtmlURL     string
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

func (c Client) Send(repoName string, attachments []slack.Attachment) {
	p := slack.Payload{
		Username:    c.username,
		IconUrl:     "",
		IconEmoji:   ":" + c.iconEmoji + ":",
		Channel:     c.channel,
		Text:        "▼ *" + repoName + "*\n",
		LinkNames:   "true",
		Attachments: attachments,
		UnfurlLinks: false,
		UnfurlMedia: false,
		Markdown:    false,
	}

	slack.Send(c.webHook, "", p)
}

func CreateMessage(prs []PullRequest, userMap []string) []slack.Attachment {
	um := make(map[string]string, 0)
	if len(userMap) > 0 {
		for _, m := range userMap {
			splitUserMap := strings.SplitN(m, ":", 2)
			um[splitUserMap[0]] = splitUserMap[1]
		}
	}

	attachments := make([]slack.Attachment, 0)
	for i, pr := range prs {
		mentions := make([]string, len(pr.Reviewers))
		for i, r := range pr.Reviewers {
			if len(userMap) > 0 {
				mentions[i] = "@" + um[r]
			} else {
				mentions[i] = "@" + r
			}
		}

		tmp := make([]string, 3)
		tmp[0] = pr.HtmlURL
		tmp[1] = strings.Join(mentions, ", ")
		tmp[2] = "Please review"
		color := "warning"
		if pr.IsMergeable {
			tmp[2] = "Let's Merge!"
			color = "good"
		}
		text := strings.Join(tmp, "\n") + "\n"

		f := slack.Field{
			Title: fmt.Sprintf("\n*%d: %s*", i+1, pr.Title),
			Value: text,
		}

		a := slack.Attachment{Color: &color}
		a.AddField(f)
		attachments = append(attachments, a)
	}
	return attachments
}
