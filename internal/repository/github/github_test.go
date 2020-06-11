package github_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/yyh-gl/gocr/internal/repository"
	"github.com/yyh-gl/gocr/internal/repository/github"
)

func Test_FetchPullRequestDetails_NotForEnterprise(t *testing.T) {
	testCases := []struct {
		name     string
		owner    string
		repoName string
		want     repository.CodeReviewRequests
	}{
		{
			name:     "Success",
			owner:    "yyh-gl",
			repoName: "gocr",
			want: &github.PullRequests{
				github.PullRequest{
					Title:          "Update go version",
					LinkURL:        "https://github.com/yyh-gl/gocr/pull/2",
					Reviewers:      []github.Reviewer{},
					Head:           github.Head{Label: "yyh-gl:samle2"},
					MergeableState: "clean",
				},
			},
		},
	}

	r, _ := recorder.New("../../../test/cassettes/FetchPullRequestDetails_NotForEnterprise")
	defer func() { _ = r.Stop() }()
	hc := http.DefaultClient
	hc.Transport = r

	repo := map[interface{}]interface{}{
		"name":  "gocr",
		"owner": "yyh-gl",
		// You set access token when recording response.
		// And you must change access token to "test-token" when commit.
		// Please also change the access token on the cassette.
		"access_token": "test-token",
		"sender":       "test",
	}
	c := github.NewClient(hc, repo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := c.FetchCodeReviewRequests(context.Background())
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
