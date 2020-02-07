package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()
	token := &oauth2.Token{AccessToken: os.Getenv("INPUT_TOKEN")}
	ts := oauth2.StaticTokenSource(token)
	tc := oauth2.NewClient(ctx, ts)
	client := &wrapper{github.NewClient(tc)}

	minDeletionSize, _ := strconv.Atoi(os.Getenv("INPUT_MINIMUMDELETIONSIZE"))
	minAge, _ := strconv.ParseFloat(os.Getenv("INPUT_MINIMUMAGE"), 64)
	name := os.Getenv("INPUT_NAME")

	ownerRepo := os.Getenv("INPUT_REPOSITORY")
	if len(ownerRepo) == 0 {
		ownerRepo = os.Getenv("GITHUB_REPOSITORY")
	}

	split := strings.SplitN(ownerRepo, "/", 2)
	owner := split[0]
	repo := split[1]

	err := forEachArtifact(ctx, client, owner, repo, func(ctx context.Context, artifact *Artifact, run *WorkflowRun) (bool, error) {
		if artifact.SizeInBytes < minDeletionSize {
			return false, nil
		}

		age := time.Now().Sub(artifact.CreatedAt)
		if age.Seconds() < minAge {
			return false, nil
		}

		if len(name) > 0 && name != artifact.Name {
			return false, nil
		}

		fmt.Printf("Deleting %s\n", artifact.URL)
		resp, err := client.DeleteWorkflowArtifact(ctx, artifact.URL)
		if err != nil {
			return true, err
		}

		if resp.StatusCode != 204 {
			return true, errors.New(fmt.Sprintf("Unexpected status code deleting artifact: %d", resp.StatusCode))
		}

		return false, nil
	})
	if err != nil {
		panic(err)
	}
}

func forEachArtifact(ctx context.Context, client *wrapper, owner, repo string, iter func(ctx context.Context, artifact *Artifact, run *WorkflowRun) (bool, error)) error {
	opt := &github.ListOptions{}

	for {
		runs, resp, err := client.ListWorkflowRuns(ctx, owner, repo, opt)
		if err != nil {
			return err
		}

		for _, run := range runs {
			artifacts, _, err := client.ListWorkflowArtifacts(ctx, run.ArtifactsURL, nil)
			if err != nil {
				return err
			}

			for _, artifact := range artifacts {
				stop, err := iter(ctx, artifact, run)
				if err != nil {
					return err
				}

				if stop {
					return nil
				}
			}
		}

		if resp.NextPage == 0 {
			return nil
		}

		opt.Page = resp.NextPage
	}
}
