package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	Client struct {
		host         string
		isEnterprise bool
		accessToken  string
	}

	PullRequest struct {
		Title          string     `json:"title"`
		URL            string     `json:"url"`
		Reviewers      []Reviewer `json:"requested_reviewers"`
		MergeableState string     `json:"mergeable_state"`
	}

	Reviewer struct {
		Login string `json:"login"`
	}
)

func NewGeneralClient(at string) *Client {
	return &Client{
		host:         "https://api.github.com",
		isEnterprise: false,
		accessToken:  at,
	}
}

func NewEnterpriseClient(h, at string) *Client {
	return &Client{
		host:         h,
		isEnterprise: true,
		accessToken:  at,
	}
}

func (c Client) FetchPullRequestDetails(owner, repo string) ([]*PullRequest, error) {
	urls, err := c.fetchOpenedPullRequestURLs(owner, repo)
	if err != nil {
		return nil, err
	}

	prs := make([]*PullRequest, 0)
	for _, u := range urls {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.accessToken)

		client := http.Client{}
		resp, err := client.Do(req)
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

		prs = append(prs, &pr)
	}

	return prs, nil
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

	client := http.Client{}
	resp, err := client.Do(req)
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
