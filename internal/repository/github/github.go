package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yyh-gl/gocr/internal/repository"
	"github.com/yyh-gl/gocr/internal/sender"
)

type (
	Repository struct {
		httpClient   *http.Client
		name         string
		owner        string
		accessToken  string
		senderID     string
		isEnterprise bool
		host         string
	}

	PullRequest struct {
		Title          string         `json:"title"`
		LinkURL        string         `json:"html_url"`
		Reviewers      []Reviewer     `json:"requested_reviewers"`
		MergeableState MergeableState `json:"mergeable_state"`
	}

	MergeableState string

	PullRequests []PullRequest

	Reviewer struct {
		Login string `json:"login"`
	}
)

func NewClient(hc *http.Client, repo interface{}) repository.Repository {
	// TODO: yaml読み取りデータのバリデーションチェック
	r := repo.(map[interface{}]interface{})

	ie, ok := r["is_enterprise"]
	isEnterprise := ok && ie.(bool)

	h := "https://api.github.com"
	if isEnterprise {
		h = r["enterprise_host"].(string)
	}
	return &Repository{
		httpClient:   hc,
		name:         r["name"].(string),
		owner:        r["owner"].(string),
		accessToken:  r["access_token"].(string),
		senderID:     r["sender"].(string),
		isEnterprise: isEnterprise,
		host:         h,
	}
}

func (r Repository) Name() string {
	return r.name
}

func (r Repository) SenderID() string {
	return r.senderID
}

func (r Repository) FetchCodeReviewRequests(ctx context.Context) (repository.CodeReviewRequests, error) {
	urls, err := r.fetchOpenedPullRequestURLs()
	if err != nil {
		return nil, err
	}

	prs := make(PullRequests, 0)
	for _, u := range urls {
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", r.accessToken)

		resp, err := r.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var pr PullRequest
		if err := json.Unmarshal(b, &pr); err != nil {
			return nil, err
		}

		prs = append(prs, pr)
	}

	return &prs, nil
}

func (r Repository) fetchOpenedPullRequestURLs() ([]string, error) {
	const endpointTemplate = "%s/repos/%s/%s/pulls?status=open"

	type pullRequest struct {
		URL string `json:"url"`
	}

	ep := fmt.Sprintf(endpointTemplate, r.host, r.owner, r.name)
	req, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", r.accessToken)

	resp, err := r.httpClient.Do(req)
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

func (pr PullRequest) ConvertToMsgMaterial() (*sender.Material, error) {
	reviwers := make([]string, len(pr.Reviewers))
	for i, r := range pr.Reviewers {
		reviwers[i] = r.Login
	}

	return &sender.Material{
		Title:       pr.Title,
		LinkURL:     pr.LinkURL,
		Reviewers:   reviwers,
		IsMergeable: pr.MergeableState.isMergeable(),
	}, nil
}

func (prs PullRequests) ConvertToMsgMaterials() (sender.Materials, error) {
	msgs := make(sender.Materials, len(prs))
	for i, pr := range prs {
		msg, err := pr.ConvertToMsgMaterial()
		if err != nil {
			return nil, err
		}
		msgs[i] = *msg
	}
	return msgs, nil
}
