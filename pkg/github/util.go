package igithub

import (
	"github.com/google/go-github/github"
)

func findFileInCommit(files []string, filename string) int {
	for i, name := range files {
		if filename == name {
			return i
		}
	}
	return len(files)
}

func isFileChangedByCommits(filename string, commits []github.PushEventCommit) (github.PushEventCommit, bool) {
	var result int
	for _, commit := range commits {
		result = findFileInCommit(commit.Added, filename)
		if result != len(commit.Added) {
			return commit, true
		}

		result = findFileInCommit(commit.Removed, filename)
		if result != len(commit.Removed) {
			return commit, true
		}

		result = findFileInCommit(commit.Modified, filename)
		if result != len(commit.Modified) {
			return commit, true
		}
	}
	return github.PushEventCommit{}, false
}
