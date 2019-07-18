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

// NewSyncWalkPath will walk the path, and store the key and full path
// of the object to be uploaded. This will return a new SyncFolderIterator
// with the data provided from walking the path.
func NewSyncWalkPath(fileconfig map[string]interface{}, bucket string) *SyncFolderIterator {
	fpath, initialrun := fileconfig[CONF_DIR_PATH].(string), fileconfig[CONF_INIT_RUN].(bool)
	metadata := []fileInfo{}
	tNow := time.Now()
	filepath.Walk(fpath, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := strings.TrimPrefix(strings.TrimPrefix(p, fpath), PathSeparator)
			if initialrun || (!initialrun && DateEqual(tNow, info.ModTime())) {
				metadata = append(metadata, fileInfo{key, p})
			} else {
				metadata = append(metadata, fileInfo{key, p})
			}
		}

		return nil
	})

	return &SyncFolderIterator{
		bucket,
		metadata,
		nil,
	}
}

func ArchiveProcess(fpath, key, targetDir, ext string) (target, fname string, e error) {
	target = fmt.Sprintf("%s.%s", filepath.Join(targetDir, key), strings.TrimLeft(ext, "."))
	fname = fmt.Sprintf("%s.%s", key, strings.TrimLeft(ext, "."))
	switch ext {
	case "zip":
		e = ZipCompress(fpath, target)
	case "tar":
		e = TarCompress(fpath, target)
	case "gz", "tar.gz":
		tarTarget := fmt.Sprintf("%s.%s", filepath.Join(targetDir, key), "tar")
		e = TarCompress(fpath, tarTarget)
		if ext == "gz" {
			target = fmt.Sprintf("%s.%s", filepath.Join(targetDir, key), "tar.gz")
		}
		e = GzCompress(tarTarget, target)
	}
	if e != nil {
		return
	}
	return
}

func NewSyncFolderIter(fileconfig map[string]interface{}, bucket string) (iter *SyncFolderIterator, tempDir string) {
	fpath, initialrun := fileconfig[CONF_DIR_PATH].(string), fileconfig[CONF_INIT_RUN].(bool)
	tempDir = filepath.Join(fpath, ArchiveTempDir)
	os.MkdirAll(tempDir, 0777)
	metadata := []fileInfo{}
	tNow := time.Now()
	filepath.Walk(fpath, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() && p != fpath && p != tempDir {
			key := strings.TrimPrefix(strings.TrimPrefix(p, fpath), PathSeparator)
			if initialrun || (!initialrun && DateEqual(tNow, info.ModTime())) {
				arcDir, arcFile, e := ArchiveProcess(p, key, tempDir, fileconfig[CONF_ARC_METHOD].(string))
				if e != nil {
					log.Println(e.Error())
				}
				metadata = append(metadata, fileInfo{arcFile, arcDir})
			}
		}

		return nil
	})
	time.Sleep(time.Second * 3) // just estimated time to wait till archiving finished
	iter = &SyncFolderIterator{
		bucket,
		metadata,
		nil,
	}

	return
}

func NewSyncFilesIter(fpath, bucket string) *SyncFolderIterator {
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
