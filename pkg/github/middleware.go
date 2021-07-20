package igithub

import (
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
	"github.com/justinas/alice"
)

type GithubMiddleware struct {
	WebEvent interface{}
	Payload  interface{}
	conf     *GithubConfig
}

func NewGithubMiddleware(conf *GithubConfig) *GithubMiddleware {
	middleware := GithubMiddleware{
		conf: conf,
	}
	return &middleware
}

func (mid *GithubMiddleware) validatePayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, err := github.ValidatePayload(r, []byte(mid.conf.GithubApp.GithubWebhookSecret))
		if err != nil {
			log.Printf("could not validate payload: err=%s\n", err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		defer r.Body.Close()

		event, err := github.ParseWebHook(github.WebHookType(r), p)
		if err != nil {
			log.Printf("could not parse webhook: err=%s\n", err)
			return
		}

		mid.WebEvent = event
		next.ServeHTTP(w, r)
	})
}

func (mid *GithubMiddleware) NewGithubInterface() alice.Chain {
	return alice.New(mid.validatePayload, mid.authenticate)
}

func NewGitHubInstallation(installationId int, conf *GithubConfig) (*github.Client, error) {
	transport := http.DefaultTransport
	itr, err := ghinstallation.New(
		transport,
		conf.GithubApp.GithubAppIdentifier,
		int64(installationId),
		[]byte(conf.GithubApp.GithubPrivateKey),
	)
	if err != nil {
		return nil, err
	}

	client := github.NewClient(&http.Client{Transport: itr})
	return client, nil
}
