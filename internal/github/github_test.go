package github_test

import (
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"

	"github.com/stretchr/testify/assert"
	"github.com/yyh-gl/gocr/internal/github"
)

func Test_GeneralClient_FetchPullRequestDetails(t *testing.T) {
	testCases := []struct {
		name     string
		owner    string
		repoName string
		want     *github.PullRequests
	}{
		{
			name:     "Success",
			owner:    "yyh-gl",
			repoName: "gocr",
			want: &github.PullRequests{
				github.PullRequest{
					Title:          "test",
					HtmlURL:        "https://github.com/yyh-gl/gocr/pull/45",
					Reviewers:      []github.Reviewer{{Login: "yyh-gl-robo"}},
					MergeableState: "clean",
				},
			},
		},
	}

	r, _ := recorder.New("../../test/cassettes/GeneralClient_FetchPullRequestDetails")
	defer r.Stop()
	hc := http.DefaultClient
	hc.Transport = r

	// You set access token when recording response.
	// And you must remove access token when commit.
	at := ""
	c := github.NewGeneralClient(hc, at)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := c.FetchPullRequestDetails(tc.owner, tc.repoName)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
