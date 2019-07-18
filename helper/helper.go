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

func BackupDBToS3(dbconfig, s3config map[string]interface{}) {
	archiveName, fPath := BackupDB(dbconfig)
	retentionDay := dbconfig["retentionday"].(float64) + 1
	if retentionDay > 0 {
		RetentionCheck(dbconfig, s3config, retentionDay)
	}
	PutObjectWithContext(s3config, archiveName, fPath)
}

func BackupDB(dbconfig map[string]interface{}) (archiveName, fPath string) {
	tNow := time.Now()
	archiveName = GenerateArchiveName(dbconfig, tNow)
	fPath = filepath.Join(dbconfig["destpath"].(string), archiveName)
	archiveCmd := "--archive=" + fPath
	ExecCommand([]string{"/C", "mongodump", "--uri", dbconfig["uri"].(string), archiveCmd})
	return
}

func GetBucketPathFromConfig(s3config map[string]interface{}) (bucket string) {
	bucket = strings.Trim(s3config["bucket"].(string), "/")
	folder := strings.Trim(s3config["folder"].(string), "/")
	if folder != "" {
		bucket = fmt.Sprintf("%s/%s/", bucket, folder)
	}
	return
}

func GetFolderPathFromConfig(s3config map[string]interface{}) string {
	return strings.Trim(s3config["folder"].(string), "/") + "/"
}

func GenerateArchiveName(dbconfig map[string]interface{}, t time.Time) string {
	return fmt.Sprintf("%s_%s.archive", dbconfig["archivename"].(string), t.Format(dbconfig["archivesuffix_dateformat"].(string)))
}
