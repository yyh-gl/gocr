package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yyh-gl/gocr/internal/slack"
)

type (
	Client struct {
		httpClient   *http.Client
		host         string
		isEnterprise bool
		accessToken  string
	}

	PullRequests []PullRequest

	PullRequest struct {
		Title          string         `json:"title"`
		HTMLURL        string         `json:"html_url"`
		Reviewers      []Reviewer     `json:"requested_reviewers"`
		MergeableState MergeableState `json:"mergeable_state"`
	}

	Reviewer struct {
		Login string `json:"login"`
	}

	MergeableState string
)

func NewGeneralClient(hc *http.Client, at string) *Client {
	return &Client{
		httpClient:   hc,
		host:         "https://api.github.com",
		isEnterprise: false,
		accessToken:  "token " + at,
	}
}

func NewEnterpriseClient(hc *http.Client, h, at string) *Client {
	return &Client{
		httpClient:   hc,
		host:         h,
		isEnterprise: true,
		accessToken:  "token " + at,
	}
}

func (c Client) FetchPullRequestDetails(owner, repo string) (*PullRequests, error) {
	urls, err := c.fetchOpenedPullRequestURLs(owner, repo)
	if err != nil {
		return nil, err
	}

	prs := make(PullRequests, 0)
	for _, u := range urls {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.accessToken)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		pr := PullRequest{}
		if err := json.Unmarshal(b, &pr); err != nil {
			return nil, err
		}

		prs = append(prs, pr)
	}

	return &prs, nil
}

func (c Client) fetchOpenedPullRequestURLs(owner, repo string) ([]string, error) {
	const endpointTemplate = "%s/repos/%s/%s/pulls?status=open"

	type pullRequest struct {
		URL string `json:"url"`
	}

	ep := fmt.Sprintf(endpointTemplate, c.host, owner, repo)
	req, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	prs := make([]pullRequest, 0)
	if err := json.Unmarshal(b, &prs); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	for _, pr := range prs {
		urls = append(urls, pr.URL)
	}
	return urls, nil
}

func (ms MergeableState) isMergeable() bool {
	return ms == "clean"
}

func (pr PullRequest) ConvertToSlackDTO() *slack.PullRequest {
	reviwers := make([]string, len(pr.Reviewers))
	for i, r := range pr.Reviewers {
		reviwers[i] = r.Login
	}

	return &slack.PullRequest{
		Title:       pr.Title,
		HTMLURL:     pr.HTMLURL,
		Reviewers:   reviwers,
		IsMergeable: pr.MergeableState.isMergeable(),
	}
}

func (prs PullRequests) ConvertToSlackDTOs() (sprs []slack.PullRequest) {
	for _, pr := range prs {
		sprs = append(sprs, *pr.ConvertToSlackDTO())
	}
	return sprs
}
