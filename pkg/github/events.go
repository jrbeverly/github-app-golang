package igithub

import (
	"context"
	b64 "encoding/base64"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/jrbeverly/github-app-golang/lib/cobrago"
)

const managedGitConfigurationFile = ".aws/config"

func (mid *GithubMiddleware) processIssueCommentEvent(event *github.IssueCommentEvent) (cobrago.TestTriggerEvent, error) {
	return cobrago.TestTriggerEvent{Key: *event.Comment.Body}, nil
}

func (mid *GithubMiddleware) processPushEvent(event *github.PushEvent) (cobrago.ConfigChangeEvent, bool, error) {
	var result cobrago.ConfigChangeEvent
	commit, ok := isFileChangedByCommits(managedGitConfigurationFile, event.Commits)
	if !ok {
		return result, false, nil
	}

	client, err := NewGitHubInstallation(int(*event.Installation.ID), mid.conf)
	if err != nil {
		return result, false, err
	}

	data, _, _, err := client.Repositories.GetContents(
		context.Background(),
		*event.Repo.Owner.Name,
		*event.Repo.Name,
		managedGitConfigurationFile,
		&github.RepositoryContentGetOptions{
			Ref: *commit.ID,
		},
	)
	if err != nil {
		return result, false, err
	}

	decodedContent, err := b64.StdEncoding.DecodeString(*data.Content)
	if err != nil {
		return result, false, err
	}

	result = cobrago.ConfigChangeEvent{
		Key: string(decodedContent),
	}
	return result, true, nil
}

func (mid *GithubMiddleware) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		event := mid.WebEvent
		switch e := event.(type) {
		case *github.IssueCommentEvent:
			result, err := mid.processIssueCommentEvent(e)
			if err != nil {
				return
			}
			mid.Payload = result
		case *github.PushEvent:
			result, ok, err := mid.processPushEvent(e)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			if !ok {
				http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
				return
			}
			mid.Payload = result
		}
		next.ServeHTTP(w, r)
	})
}
