package objectstorage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ObjectStorage interface {
	UploadMediaFiles(prefix, localPath string) error
	DeleteMediaFiles(prefix string) error
}

type objectStorage struct {
	sess   *session.Session
	bucket string
}

func NewObjectStorage(accessKey, secretKey, endpoint, region, bucket string) (ObjectStorage, error) {
	bucketSession, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
	if err != nil {
		return nil, err
	}
	return &objectStorage{
		sess:   bucketSession,
		bucket: bucket,
	}, nil
}

func (o *objectStorage) UploadMediaFiles(prefix, localPath string) error {
	client := s3.New(o.sess)
	log.Println("Removing existing files on the bucket on path", prefix)
	err := o.deleteDirectoryFromS3(client, prefix)
	if err != nil {
		return err
	}
	log.Println("Uploading files to the bucket on path", prefix)
	err = o.uploadDirectoryToS3(client, prefix, localPath)
	if err != nil {
		return err
	}
	log.Println("Files uploaded successfully")
	return nil
}

func (o *objectStorage) DeleteMediaFiles(prefix string) error {
	client := s3.New(o.sess)
	log.Println("Removing existing files on the bucket on path", prefix)
	err := o.deleteDirectoryFromS3(client, prefix)
	if err != nil {
		return err
	}
	log.Println("Files removed successfully")
	return nil
}
func (o *objectStorage) deleteDirectoryFromS3(client *s3.S3, prefix string) error {
	var continuationToken *string

	for {
		objects, token, err := o.listObjectsForDeletion(client, prefix, continuationToken)
		if err != nil {
			return err
		}

		if len(objects) > 0 {
			err = o.deleteObjects(client, objects)
			if err != nil {
				return err
			}
		}

		if token == nil {
			break
		} else {
			continuationToken = token
		}
	}

	return nil
}

func (o *objectStorage) listObjectsForDeletion(client *s3.S3, prefix string, continuationToken *string) ([]*s3.ObjectIdentifier, *string, error) {
	resp, err := client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:            aws.String(o.bucket),
		Prefix:            aws.String(prefix),
		ContinuationToken: continuationToken,
	})
	if err != nil {
		log.Println("Error while listing objects for deletion in", prefix)
		return nil, nil, err
	}

	objects := make([]*s3.ObjectIdentifier, len(resp.Contents))
	for i, item := range resp.Contents {
		objects[i] = &s3.ObjectIdentifier{Key: item.Key}
	}

	return objects, resp.NextContinuationToken, nil
}

func (o *objectStorage) deleteObjects(client *s3.S3, objects []*s3.ObjectIdentifier) error {
	var err error
	for i := 0; i < 3; i++ {
		_, err = client.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(o.bucket),
			Delete: &s3.Delete{
				Objects: objects,
				Quiet:   aws.Bool(true),
			},
		})

		if err != nil {
			log.Printf("Error while removing objects, error: %s\nRetrying...", err.Error())
			time.Sleep(1 * time.Second) // wait for 1 second before next attempt
		} else {
			return nil
		}
	}

	log.Println("Failed to delete objects after 3 attempts")
	return err
}

func (o *objectStorage) uploadFileToS3(client *s3.S3, prefix, filePath string, wg *sync.WaitGroup, sem chan bool) {
	defer wg.Done()

	sem <- true // block until there's room
	defer func() { <-sem }()

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Failed to open file", filePath)
		return
	}
	defer file.Close()

	_, filename := filepath.Split(filePath)
	key := filepath.Join(prefix, filename)

	var success bool
	for i := 0; i < 3; i++ {
		_, err = client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(o.bucket),
			Key:    aws.String(key),
			ACL:    aws.String("public-read"),
			Body:   file,
		})

		if err != nil {
			log.Printf("Failed to upload %s to bucket %s, error: %s\nRetrying...", key, o.bucket, err.Error())
			time.Sleep(1 * time.Second) // wait for 1 second before next attempt
		} else {
			//log.Printf("File %s uploaded successfully", key)
			success = true
			break
		}
	}

	if !success {
		log.Printf("Failed to upload %s to bucket %s after 3 attempts", key, o.bucket)
	}
}

func (o *objectStorage) uploadDirectoryToS3(client *s3.S3, prefix, localPath string) error {
	var wg sync.WaitGroup
	sem := make(chan bool, 4) // limit to 4 concurrent goroutines
	log.Println("Uploading files from", localPath, "to", prefix)
	err := filepath.WalkDir(localPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			wg.Add(1)
			go func(path string) {
				o.uploadFileToS3(client, prefix, path, &wg, sem)
			}(path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory %s, error: %s", localPath, err.Error())
	}

	wg.Wait()

	return nil
}
