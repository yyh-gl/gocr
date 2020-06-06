package slack

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/yyh-gl/gocr/internal/sender"
)

type (
	Sender struct {
		webHook   string
		channel   string
		username  string
		iconEmoji string
		userMap   []interface{}
	}
)

func NewClient(sender interface{}) sender.Sender {
	// TODO: yaml読み取りデータのバリデーションチェック
	s := sender.(map[interface{}]interface{})

	// TODO: refactoring
	var userMap []interface{}
	if _, ok := s["user_map"]; ok {
		userMap = s["user_map"].([]interface{})
	}

	return &Sender{
		webHook:   s["web_hook"].(string),
		channel:   s["channel"].(string),
		username:  s["username"].(string),
		iconEmoji: s["icon_emoji"].(string),
		userMap:   userMap,
	}
}

func (s Sender) Send(repoName string, materials sender.Materials) error {
	attachments := createMessage(materials, s.userMap)
	p := slack.Payload{
		Username:    s.username,
		IconUrl:     "",
		IconEmoji:   ":" + s.iconEmoji + ":",
		Channel:     s.channel,
		Text:        "▼ *" + repoName + "*\n",
		LinkNames:   "true",
		Attachments: attachments,
		UnfurlLinks: false,
		UnfurlMedia: false,
		Markdown:    false,
	}

	slack.Send(s.webHook, "", p)
	return nil
}

func createMessage(materials sender.Materials, userMap []interface{}) []slack.Attachment {
	um := make(map[string]string, 0)
	if len(userMap) > 0 {
		for _, m := range userMap {
			splitUserMap := strings.SplitN(m.(string), ":", 2)
			um[splitUserMap[0]] = splitUserMap[1]
		}
	}

	attachments := make([]slack.Attachment, 0)
	for i, m := range materials {
		mentionStr := "no reviewer"

		// TODO: refactoring. make more simple
		reviewersCount := len(m.Reviewers)
		if reviewersCount > 0 {
			mentions := make([]string, reviewersCount)
			for i, r := range m.Reviewers {
				if len(userMap) > 0 {
					mentions[i] = "@" + um[r]
				} else {
					mentions[i] = "@" + r
				}
			}
			mentionStr = strings.Join(mentions, ", ")
		}

		tmp := make([]string, 3)
		tmp[0] = m.LinkURL
		tmp[1] = mentionStr
		tmp[2] = "Please review"
		color := "warning"
		if m.IsMergeable {
			tmp[2] = "Let's Merge!"
			color = "good"
		}
		text := strings.Join(tmp, "\n") + "\n"

		f := slack.Field{
			Title: fmt.Sprintf("\n*%d: %s*", i+1, m.Title),
			Value: text,
		}

		a := slack.Attachment{Color: &color}
		a.AddField(f)
		attachments = append(attachments, a)
	}
	return attachments
}
