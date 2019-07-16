package helper

import (
	"fmt"
	"path/filepath"
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
	// if dbconfig["retentionday"].(float64) > 0 {
	// 	GetListObjects(s3config)
	// }
	PutObjectWithContext(s3config, archiveName, fPath)
}

func BackupDB(dbconfig map[string]interface{}) (archiveName, fPath string) {
	tNow := time.Now()
	archiveName = fmt.Sprintf("%s_%s.archive", dbconfig["archivename"].(string), tNow.Format(dbconfig["archivesuffix_dateformat"].(string)))
	fPath = filepath.Join(dbconfig["destpath"].(string), archiveName)
	archiveCmd := "--archive=" + fPath
	ExecCommand([]string{"/C", "mongodump", "--uri", dbconfig["uri"].(string), archiveCmd})
	return
}
