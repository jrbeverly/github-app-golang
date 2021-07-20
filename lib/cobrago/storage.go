package cobrago

import "fmt"

type RemoteStorage interface {
	List(bucket string) []RemoteFile
}

type SystemWriter interface {
	PrintRemoteFiles(files []RemoteFile)
	PrintAWSConfiguration(config ConfigChangeEvent)
	PrintTestResults(results TestResults)
}

type RemoteFile struct {
	Key  string
	Size int64
}

func ListFilesFromStorage(bucket string, storage RemoteStorage, writer SystemWriter) {
	files := storage.List(bucket)
	for _, file := range files {
		file.Key = fmt.Sprintf("Prefix: %s", file.Key)
	}
	writer.PrintRemoteFiles(files)
}

type TestTriggerEvent struct {
	Key string
}

type ConfigChangeEvent struct {
	Key string
}

type TestResults struct {
	Success bool
}

func PerformTestTrigger(trigger TestTriggerEvent, writer SystemWriter) {
	writer.PrintTestResults(TestResults{Success: true})
}

func PerformConfigTrigger(trigger ConfigChangeEvent, writer SystemWriter) {
	writer.PrintAWSConfiguration(trigger)
}
