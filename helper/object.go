package helper

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
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

const (
	ApEast1RegionID      = "ap-east-1"      // Asia Pacific (Hong Kong).
	ApNortheast1RegionID = "ap-northeast-1" // Asia Pacific (Tokyo).
	ApNortheast2RegionID = "ap-northeast-2" // Asia Pacific (Seoul).
	ApSouth1RegionID     = "ap-south-1"     // Asia Pacific (Mumbai).
	ApSoutheast1RegionID = "ap-southeast-1" // Asia Pacific (Singapore).
	ApSoutheast2RegionID = "ap-southeast-2" // Asia Pacific (Sydney).
	CaCentral1RegionID   = "ca-central-1"   // Canada (Central).
	EuCentral1RegionID   = "eu-central-1"   // EU (Frankfurt).
	EuNorth1RegionID     = "eu-north-1"     // EU (Stockholm).
	EuWest1RegionID      = "eu-west-1"      // EU (Ireland).
	EuWest2RegionID      = "eu-west-2"      // EU (London).
	EuWest3RegionID      = "eu-west-3"      // EU (Paris).
	SaEast1RegionID      = "sa-east-1"      // South America (Sao Paulo).
	UsEast1RegionID      = "us-east-1"      // US East (N. Virginia).
	UsEast2RegionID      = "us-east-2"      // US East (Ohio).
	UsWest1RegionID      = "us-west-1"      // US West (N. California).
	UsWest2RegionID      = "us-west-2"      // US West (Oregon).
)

func PutObjectAcl(regionPtr, bucketPtr, keyPtr, ownerNamePtr, ownerIDPtr, granteeTypePtr, uriPtr, emailPtr, userPtr, displayNamePtr *string) {
	// Based off the type, fields must be excluded.
	switch *granteeTypePtr {
	case s3.TypeCanonicalUser:
		emailPtr, uriPtr = nil, nil
		if *displayNamePtr == "" {
			displayNamePtr = nil
		}

		if *userPtr == "" {
			userPtr = nil
		}
	case s3.TypeAmazonCustomerByEmail:
		uriPtr, userPtr = nil, nil
	case s3.TypeGroup:
		emailPtr, userPtr = nil, nil
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: regionPtr,
	}))

	svc := s3.New(sess)

	resp, err := svc.PutObjectAcl(&s3.PutObjectAclInput{
		Bucket: bucketPtr,
		Key:    keyPtr,
		AccessControlPolicy: &s3.AccessControlPolicy{
			Owner: &s3.Owner{
				DisplayName: ownerNamePtr,
				ID:          ownerIDPtr,
			},
			Grants: []*s3.Grant{
				{
					Grantee: &s3.Grantee{
						Type:         granteeTypePtr,
						DisplayName:  displayNamePtr,
						URI:          uriPtr,
						EmailAddress: emailPtr,
						ID:           userPtr,
					},
					Permission: aws.String(s3.BucketLogsPermissionFullControl),
				},
			},
		},
	})

	if err != nil {
		log.Println("failed", err)
	} else {
		log.Println("success", resp)
	}
}

func GetListObjects(s3config map[string]interface{}) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config["region"].(string)),
	}))
	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(s3config["bucket"].(string)),
		MaxKeys: aws.Int64(2),
	}

	result, e := svc.ListObjectsV2(input)
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
	for _, obj := range result.Contents {
		log.Println(*obj.Key, obj.LastModified)
	}
}

func PutObjectWithContext(s3config map[string]interface{}, key, fPath string) {
	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s3config["region"].(string)),
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
	timeout := time.Duration(s3config["timeout"].(float64))
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

	bucket := strings.Trim(s3config["bucket"].(string), "/")
	folder := strings.Trim(s3config["folder"].(string), "/")
	if folder != "" {
		bucket = fmt.Sprintf("%s/%s/", bucket, folder)
	}
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

	log.Printf("successfully uploaded file to %s-%s\n", bucket, key)
}

func PutObjectsToS3(fileconfig, s3config map[string]interface{}) {
	sess := session.New(&aws.Config{
		Region: aws.String(s3config["region"].(string)),
	})
	uploader := s3manager.NewUploader(sess)

	bucket := strings.Trim(s3config["bucket"].(string), "/")
	folder := strings.Trim(s3config["folder"].(string), "/")
	if folder != "" {
		bucket = fmt.Sprintf("%s/%s/", bucket, folder)
	}
	iter := NewSyncFolderIter(fileconfig["dirpath"].(string), bucket)
	if err := uploader.UploadWithIterator(aws.BackgroundContext(), iter); err != nil {
		log.Printf("unexpected error has occurred: %v", err)
	}

	if err := iter.Err(); err != nil {
		log.Printf("unexpected error occurred during file walking: %v", err)
	}

	log.Println("Backup File to S3 success")
}

// SyncFolderIterator is used to upload a given folder
// to Amazon S3.
type SyncFolderIterator struct {
	bucket    string
	fileInfos []fileInfo
	err       error
}

type fileInfo struct {
	key      string
	fullpath string
}

func NewSyncFolderIter(fpath, bucket string) *SyncFolderIterator {
	metadata := []fileInfo{}
	tNow := time.Now()
	files, e := ioutil.ReadDir(fpath)
	if e != nil {
		log.Println(e.Error())
		return nil
	}

	for _, f := range files {
		if DateEqual(tNow, f.ModTime()) {
			metadata = append(metadata, fileInfo{f.Name(), filepath.Join(fpath, f.Name())})
		}
	}

	return &SyncFolderIterator{
		bucket,
		metadata,
		nil,
	}
}

// Next will determine whether or not there is any remaining files to
// be uploaded.
func (iter *SyncFolderIterator) Next() bool {
	return len(iter.fileInfos) > 0
}

// Err returns any error when os.Open is called.
func (iter *SyncFolderIterator) Err() error {
	return iter.err
}

// UploadObject will prep the new upload object by open that file and constructing a new
// s3manager.UploadInput.
func (iter *SyncFolderIterator) UploadObject() s3manager.BatchUploadObject {
	fi := iter.fileInfos[0]
	iter.fileInfos = iter.fileInfos[1:]
	body, err := os.Open(fi.fullpath)
	if err != nil {
		iter.err = err
	}

	extension := filepath.Ext(fi.key)
	mimeType := mime.TypeByExtension(extension)

	if mimeType == "" {
		mimeType = "binary/octet-stream"
	}

	input := s3manager.UploadInput{
		Bucket:      &iter.bucket,
		Key:         &fi.key,
		Body:        body,
		ContentType: &mimeType,
	}

	return s3manager.BatchUploadObject{
		Object: &input,
	}
}

// NewSyncFolderIterator will walk the path, and store the key and full path
// of the object to be uploaded. This will return a new SyncFolderIterator
// with the data provided from walking the path.
// func NewSyncFolderIterator(path, bucket string) *SyncFolderIterator {
// 	metadata := []fileInfo{}
// 	timeNow := time.Now()
// 	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
// 		if !info.IsDir() {
// 			key := strings.TrimPrefix(p, path)
// 			metadata = append(metadata, fileInfo{key, p})
// 		}

// 		return nil
// 	})

// 	return &SyncFolderIterator{
// 		bucket,
// 		metadata,
// 		nil,
// 	}
// }
