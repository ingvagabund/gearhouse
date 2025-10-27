package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v76/github"
	"golang.org/x/oauth2"

	"k8s.io/klog/v2"
)

const (
	updateKonfluxReferencesPRTitle  = "chore(deps): update konflux references"
	updateKonfluxReferencesPRTitle2 = "Update Konflux references"
	updateDockerfileBundlePRTitle   = "chore(deps): update"
)

var allowedAuthors = map[string]bool{
	"red-hat-konflux[bot]": true,
}

var label2comments = map[string]string{
	"ok-to-test":             "/ok-to-test",
	"backport-risk-assessed": "/label backport-risk-assessed",
	"verified":               "/verified by CI",
	"lgtm":                   "/lgtm",
	"approved":               "/approved",
}

func getChangedFiles(ctx context.Context, client *github.Client, owner, repo string, prNum int) ([]string, error) {
	allFiles := []string{}
	listOpts := &github.ListOptions{PerPage: 100}

	for {
		files, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, prNum, listOpts)
		if err != nil {
			return nil, fmt.Errorf("error listing files for PR #%d: %v", prNum, err)
		}

		for _, file := range files {
			allFiles = append(allFiles, file.GetFilename())
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	return allFiles, nil
}

func validateUpdateKonfluxReferences(files []string) bool {
	for _, file := range files {
		if !strings.HasPrefix(file, ".tekton") {
			return false
		}
		if !strings.HasSuffix(file, ".yaml") {
			return false
		}
		if strings.HasSuffix(file, "images-mirror-set.yaml") {
			return false
		}
	}
	return true
}

func validateUpdateBundleImageShas(files []string) bool {
	for _, file := range files {
		if file != "bundle.Dockerfile" {
			return false
		}
	}
	return true
}

func ensurePRLabels(ctx context.Context, client *github.Client, owner, repo string, prNum int, pr *github.PullRequest, labels []string) error {
	existingLabels := make(map[string]struct{})
	for _, label := range pr.Labels {
		existingLabels[label.GetName()] = struct{}{}
	}

	mustHaveLabels := []string{}
	for _, label := range labels {
		if _, exists := existingLabels[label]; !exists {
			mustHaveLabels = append(mustHaveLabels, label)
		}
	}

	if len(mustHaveLabels) == 0 {
		klog.InfoS("No missing labels for PR", "number", prNum)
		return nil
	}

	klog.InfoS("Adding missing labels to PR", "number", prNum, "labels", mustHaveLabels)
	_, _, err := client.Issues.AddLabelsToIssue(ctx, owner, repo, prNum, mustHaveLabels)
	if err != nil {
		return fmt.Errorf("error adding label to PR #%d: %v", prNum, err)
	}

	return nil
}

func ensurePRCommentBasedLabel(ctx context.Context, client *github.Client, owner, repo string, prNum int, pr *github.PullRequest, targetLabel, targetComment string) error {
	for _, label := range pr.Labels {
		if label.GetName() == targetLabel {
			klog.InfoS("Label already present for PR", "number", prNum, "label", targetLabel)
			return nil
		}
	}

	comment := &github.IssueComment{
		Body: github.String(targetComment),
	}

	klog.InfoS("Adding comment to PR", "number", prNum, "comment", targetComment)
	_, _, err := client.Issues.CreateComment(ctx, owner, repo, prNum, comment)
	return err
}

func inspectRepository(ctx context.Context, client *github.Client, organization, repository string) {
	klog.Infof("Fetching open Pull Requests for %s/%s...", organization, repository)

	opts := &github.PullRequestListOptions{
		State:       "open",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	prs, _, err := client.PullRequests.List(ctx, organization, repository, opts)
	if err != nil {
		log.Fatalf("Error listing PRs: %v", err)
	}

	klog.Infof("Found %d open PRs.", len(prs))

	for _, pr := range prs {
		if pr.Number == nil || pr.User == nil || pr.User.Login == nil || pr.Title == nil {
			continue
		}

		prNum := *pr.Number
		prAuthor := *pr.User.Login

		klog.InfoS("Processing PR", "number", prNum, "author", prAuthor, "title", *pr.Title)

		if !allowedAuthors[prAuthor] {
			continue
		}

		files, err := getChangedFiles(ctx, client, organization, repository, prNum)
		if err != nil {
			klog.Errorf("Error listing files: %v", err)
			continue
		}

		if len(files) == 0 {
			continue
		}

		// Only PRs either changing just .tekton files or just Dockerfiles
		if strings.Contains(*pr.Title, updateKonfluxReferencesPRTitle) || strings.Contains(*pr.Title, updateKonfluxReferencesPRTitle2) {
			if validateUpdateKonfluxReferences(files) {
				// Set the right labels
				if err := ensurePRLabels(ctx, client, organization, repository, prNum, pr, []string{"jira/valid-bug"}); err != nil {
					klog.Errorf("Error labeling PR: %v", err)
				}
				// Produce the right labels through comments
				for targetLabel, targetComment := range label2comments {
					if err := ensurePRCommentBasedLabel(ctx, client, organization, repository, prNum, pr, targetLabel, targetComment); err != nil {
						klog.Errorf("Error ensuring %q label: %v", targetLabel, err)
					}
				}
			} else {
				klog.InfoS("validateUpdateKonfluxReferences: [false]")
			}
		}
		if strings.Contains(*pr.Title, updateDockerfileBundlePRTitle) {
			if validateUpdateBundleImageShas(files) {
				// Set the right labels
				if err := ensurePRLabels(ctx, client, organization, repository, prNum, pr, []string{"jira/valid-bug"}); err != nil {
					klog.Errorf("Error labeling PR: %v", err)
				}
				// Produce the right labels through comments
				for targetLabel, targetComment := range label2comments {
					if err := ensurePRCommentBasedLabel(ctx, client, organization, repository, prNum, pr, targetLabel, targetComment); err != nil {
						klog.Errorf("Error ensuring %q label: %v", targetLabel, err)
					}
				}
			} else {
				klog.InfoS("validateUpdateBundleImageShas: [false]")
			}
		}
	}
}

func main() {
	initFlags()
	validateFlags()

	ctx := context.Background()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Error: GITHUB_TOKEN environment variable is not set. Please set your Personal Access Token.")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for _, repo := range repositories {
		items := strings.Split(repo, "/")
		klog.InfoS("Processing repository", "organization", items[0], "repository", items[1])
		inspectRepository(ctx, client, items[0], items[1])
	}

}
