package console

import (
	"log"

	"github.com/jrbeverly/github-app-golang/lib/cobrago"
)

type ConsoleWriter struct {
}

func NewConsoleWriter() ConsoleWriter {
	writer := ConsoleWriter{}
	return writer
}

func (r ConsoleWriter) PrintRemoteFiles(files []cobrago.RemoteFile) {
	log.Println("[EVENT]: List of files:")
	for _, object := range files {
		log.Printf("[EVENT]: key=%s size=%d\n", object.Key, object.Size)
	}
}

func (r ConsoleWriter) PrintAWSConfiguration(config cobrago.ConfigChangeEvent) {
	log.Printf("[EVENT]: %s\n", config.Key)
}

func (r ConsoleWriter) PrintTestResults(results cobrago.TestResults) {
	log.Printf("[EVENT]: Test results: %v\n", results.Success)
}
