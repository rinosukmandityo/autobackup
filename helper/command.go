package helper

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"
)

func ExecCommand(commands []string) {
	if e := exec.Command("cmd", commands...).Run(); e != nil {
		log.Println(e.Error())
	}
	log.Println("dump success")
}

func BackupDBToS3(dbconfig, s3config map[string]interface{}) {
	archiveName, fPath := BackupDB(dbconfig)
	PutObjectWithContext(s3config, archiveName, fPath)
}

func BackupDB(dbconfig map[string]interface{}) (archiveName, fPath string) {
	tNow := time.Now()
	archiveName = fmt.Sprintf("%s_%s.archive", dbconfig["archivename"].(string), tNow.Format("20060102"))
	fPath = filepath.Join(dbconfig["destpath"].(string), archiveName)
	archiveCmd := "--archive=" + fPath
	ExecCommand([]string{"/C", "mongodump", "--uri", dbconfig["uri"].(string), archiveCmd})
	return
}
