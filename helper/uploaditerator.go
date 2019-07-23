package helper

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

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
