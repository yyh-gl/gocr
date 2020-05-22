package prn

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

func (c Client) FetchPRs(owner, repo string) ([]*PullRequest, error) {
	const endpointTemplate = "%s/repos/%s/%s/pulls"

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

	prs := make([]*PullRequest, 5)
	if err := json.Unmarshal(b, &prs); err != nil {
		return nil, err
	}
	return prs, nil
}
