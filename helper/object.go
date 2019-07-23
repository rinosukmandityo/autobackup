package helper

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func GetListObjectsWithContext(s3config map[string]interface{}) (result *s3.ListObjectsV2Output, e error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config[CONF_REGION].(string)),
	}))
	svc := s3.New(sess)

	ctx := context.Background()

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s3config[CONF_BUCKET].(string)),
		MaxKeys: aws.Int64(2),
	}

	result, e = svc.ListObjectsV2WithContext(ctx, input)
	if e != nil {
		if aerr, ok := e.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(e.Error())
		}
		return
	}
	return
}

func RetentionCheck(dbconfig, s3config map[string]interface{}, retentionDay float64) {
	tr := time.Now().Add(time.Hour * 24 * time.Duration(retentionDay) * -1)
	tRetention := time.Date(tr.Year(), tr.Month(), tr.Day(), 0, 0, 0, 0, tr.Location())
	archiveName := GenerateArchiveName(dbconfig, tRetention)
	fPath := filepath.Join(dbconfig[CONF_DEST_PATH].(string), archiveName)
	os.RemoveAll(fPath)
	DeleteObjectWithContext(s3config, archiveName)
}

func DeleteObjectWithContext(s3config map[string]interface{}, key string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config[CONF_REGION].(string)),
	}))
	svc := s3.New(sess)
	ctx := context.Background()

	bucket := GetBucketPathFromConfig(s3config)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	_, err := svc.DeleteObjectWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return
	}

	log.Println(fmt.Sprintf("Delete object %s success", key))
}

func DeleteObjectsWithContext(s3config map[string]interface{}, keys []string) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config[CONF_REGION].(string)),
	}))
	svc := s3.New(sess)
	ctx := context.Background()

	bucket := GetBucketPathFromConfig(s3config)

	objects := []*s3.ObjectIdentifier{}
	for _, k := range keys {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(k)})
	}
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(false),
		},
	}

	_, err := svc.DeleteObjectsWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
		return
	}

	log.Println(fmt.Sprintf("Delete object %s success", strings.Join(keys, ", ")))
}

func RetentionCheckByResult(result *s3.ListObjectsV2Output) {
	tr := time.Now().Add(time.Hour * 24 * 7 * -1)
	tRetention := time.Date(tr.Year(), tr.Month(), tr.Day(), 0, 0, 0, 0, tr.Location())
	GetDeletedObjects(result, tRetention)
}

func GetDeletedObjects(result *s3.ListObjectsV2Output, tRetention time.Time) (res []string) {
	res = []string{}
	for _, obj := range result.Contents {
		key := *obj.Key
		if !strings.HasSuffix(key, "/") {
			if obj.LastModified.Before(tRetention) {
				res = append(res, key)
			}
		}
	}
	return
}

func PutObjectWithContext(s3config map[string]interface{}, key, fPath string) {
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config[CONF_REGION].(string)),
	}))

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.
	svc := s3.New(sess)

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	timeout := time.Duration(s3config[CONF_TIMEOUT].(float64))
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
		defer cancelFn()
	}

	f, e := os.Open(fPath)
	if e != nil {
		log.Println(e.Error())
	}
	defer f.Close()

	// Uploads the object to S3. The Context will interrupt the request if the
	// timeout expires.

	bucket := GetBucketPathFromConfig(s3config)
	_, e = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if e != nil {
		if aerr, ok := e.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			log.Printf("upload canceled due to timeout, %s\n", e.Error())
		} else {
			log.Printf("failed to upload object, %s\n", e.Error())
		}
		os.Exit(1)
	}

	log.Printf("successfully uploaded file to %s%s\n", bucket, key)
}

func PutObjectsToS3(fileconfig, s3config map[string]interface{}) {
	sess := session.New(&aws.Config{
		Region: aws.String(s3config[CONF_REGION].(string)),
	})
	uploader := s3manager.NewUploader(sess)

	bucket := GetBucketPathFromConfig(s3config)
	iter := new(SyncFolderIterator)
	backupType := fileconfig[CONF_BACK_TYPE].(string)
	tempDir := ""

	switch backupType {
	case CONF_FOLDER:
		iter, tempDir = NewSyncFolderIter(fileconfig, bucket)
	case "file":
		iter = NewSyncWalkPath(fileconfig, bucket)
	}

	if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		log.Printf("unexpected error has occurred: %v", err)
	}

	if err := iter.Err(); err != nil {
		log.Printf("unexpected error occurred during file walking: %v", err)
	}

	if backupType == CONF_FOLDER {
		time.Sleep(time.Second * 2)
		os.RemoveAll(tempDir)
	}
	log.Println("Backup File to S3 success")
}
