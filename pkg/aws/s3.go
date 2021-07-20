package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jrbeverly/github-app-golang/lib/cobrago"
)

type S3RemoteStorage struct {
	Configuration *aws.Config
	S3Client      *s3.Client
}

func NewS3RemoteStorage() S3RemoteStorage {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	storage := S3RemoteStorage{
		Configuration: &cfg,
		S3Client:      client,
	}
	return storage
}

func (r S3RemoteStorage) List(bucket string) []cobrago.RemoteFile {
	// A method that interacts with AWS, then returns some results

	output, err := r.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Fatal(err)
	}

	result := make([]cobrago.RemoteFile, len(output.Contents))
	for idx, object := range output.Contents {
		file := cobrago.RemoteFile{}
		file.Key = *object.Key
		file.Size = object.Size
		result[idx] = file
	}
	return result
}
