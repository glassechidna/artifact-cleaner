package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/google/go-querystring/query"
	"net/url"
	"reflect"
	"time"
)

type wrapper struct {
	*github.Client
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func (w *wrapper) ListWorkflowRuns(ctx context.Context, owner string, repo string, opt *github.ListOptions) ([]*WorkflowRun, *github.Response, error) {
	u := fmt.Sprintf("repos/%s/%s/actions/runs", owner, repo)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := w.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var runs workflowRuns
	resp, err := w.Do(ctx, req, &runs)
	return runs.WorkflowRuns, resp, err
}

func (w *wrapper) ListWorkflowArtifacts(ctx context.Context, u string, opt *github.ListOptions) ([]*Artifact, *github.Response, error) {
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := w.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var artifacts workflowRunArtifacts
	resp, err := w.Do(ctx, req, &artifacts)
	return artifacts.Artifacts, resp, err
}

func (w *wrapper) DeleteWorkflowArtifact(ctx context.Context, u string) (*github.Response, error) {
	req, err := w.NewRequest("DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return w.Do(ctx, req, nil)
}

type workflowRuns struct {
	TotalCount   int            `json:"total_count"`
	WorkflowRuns []*WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	ID             int                     `json:"id"`
	NodeID         string                  `json:"node_id"`
	HeadBranch     string                  `json:"head_branch"`
	HeadSha        string                  `json:"head_sha"`
	RunNumber      int                     `json:"run_number"`
	CheckSuiteID   int                     `json:"check_suite_id"`
	Event          string                  `json:"event"`
	Status         string                  `json:"status"`
	Conclusion     string                  `json:"conclusion"`
	URL            string                  `json:"url"`
	HTMLURL        string                  `json:"html_url"`
	PullRequests   []interface{}           `json:"pull_requests"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
	JobsURL        string                  `json:"jobs_url"`
	LogsURL        string                  `json:"logs_url"`
	ArtifactsURL   string                  `json:"artifacts_url"`
	CancelURL      string                  `json:"cancel_url"`
	RerunURL       string                  `json:"rerun_url"`
	WorkflowURL    string                  `json:"workflow_url"`
	HeadCommit     *github.PushEventCommit `json:"head_commit"`
	Repository     *github.Repository      `json:"repository"`
	HeadRepository *github.Repository      `json:"head_repository"`
}

type workflowRunArtifacts struct {
	TotalCount int        `json:"total_count"`
	Artifacts  []*Artifact `json:"artifacts"`
}

type Artifact struct {
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	SizeInBytes        int       `json:"size_in_bytes"`
	URL                string    `json:"url"`
	ArchiveDownloadURL string    `json:"archive_download_url"`
	Expired            bool      `json:"expired"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
