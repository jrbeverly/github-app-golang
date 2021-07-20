package cmd

import (
	"net/http"

	"github.com/jrbeverly/github-app-golang/lib/cobrago"
	igithub "github.com/jrbeverly/github-app-golang/pkg/github"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var githubConfigFilePath string

var (
	conf       *igithub.GithubConfig
	middleware *igithub.GithubMiddleware
)

func init() {
	log.SetFormatter(&log.TextFormatter{})

	startCmd.Flags().StringVar(&githubConfigFilePath, "github-config", "", "Path to the GitHub configuration file")
	startCmd.MarkFlagRequired("github-config")

	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			return
		})

		config, err := igithub.ReadGithubConfig(githubConfigFilePath)
		if err != nil {
			log.Fatal("Error reading config file: ", err)
		}

		middleware = igithub.NewGithubMiddleware(&config)
		stdChain := middleware.NewGithubInterface()
		http.Handle("/event_handler", stdChain.Then(http.HandlerFunc(event_handler)))

		log.Info("Server listening...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	},
}

func event_handler(w http.ResponseWriter, r *http.Request) {
	switch e := middleware.Payload.(type) {
	case cobrago.TestTriggerEvent:
		cobrago.PerformTestTrigger(e, writer)
	case cobrago.ConfigChangeEvent:
		cobrago.PerformConfigTrigger(e, writer)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
