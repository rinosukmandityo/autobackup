package helper

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

func DateEqual(date1, date2 time.Time) bool {
	return date1.Format("20060102") == date2.Format("20060102")
}

func ToMapString(data interface{}) (res map[string]string) {
	res = map[string]string{}

	for k, v := range data.(map[string]interface{}) {
		res[k] = v.(string)
	}
	return
}

// To store your database into S3
func BackupDBToS3(dbconfig, s3config map[string]interface{}) {
	archiveName, fPath := BackupDB(dbconfig)
	retentionDay := dbconfig[CONF_RETENTION_DAY].(float64)
	if retentionDay > 0 {
		RetentionCheck(dbconfig, s3config, retentionDay+1)
	}
	PutObjectWithContext(s3config, archiveName, fPath)
}

func BackupDB(dbconfig map[string]interface{}) (archiveName, fPath string) {
	tNow := time.Now()
	archiveName = GenerateArchiveName(dbconfig, tNow)
	fPath = filepath.Join(dbconfig[CONF_DEST_PATH].(string), archiveName)
	archiveCmd := "--archive=" + fPath
	ExecCommand([]string{"/C", "mongodump", "--uri", dbconfig[CONF_URI].(string), archiveCmd})
	return
}

func GetBucketPathFromConfig(s3config map[string]interface{}) (bucket string) {
	bucket = strings.Trim(s3config[CONF_BUCKET].(string), "/")
	folder := strings.Trim(s3config[CONF_FOLDER].(string), "/")
	if folder != "" {
		bucket = fmt.Sprintf("%s/%s/", bucket, folder)
	}
	return
}

func GetFolderPathFromConfig(s3config map[string]interface{}) string {
	return strings.Trim(s3config[CONF_FOLDER].(string), "/") + "/"
}

func GenerateArchiveName(dbconfig map[string]interface{}, t time.Time) string {
	return fmt.Sprintf("%s_%s.archive", dbconfig[CONF_ARC_NAME].(string), t.Format(dbconfig[CONF_DATE_FORMAT].(string)))
}
