package helper

import (
	"os"
)

var (
	PathSeparator  string = string(os.PathSeparator)
	ArchiveTempDir        = "archivetemp"
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

const (
	CONF_FOR_FILE = "file"
	CONF_FOR_DB   = "database"
	CONF_FOR_S3   = "s3"

	CONF_ARC_METHOD = "archivemethod" // archive method zip, tar or tar.gz
	CONF_BACK_TYPE  = "backuptype"    // backup type `folder` or `file`
	CONF_DIR_PATH   = "dirpath"       // directory path of your files
	CONF_INIT_RUN   = "initialrun"    // initial run true or false

	CONF_URI           = "uri"                      // connection string of your mongoDB
	CONF_ARC_NAME      = "archivename"              // archive name of your database dump files directory
	CONF_DATE_FORMAT   = "archivesuffix_dateformat" // date format for arhive name suffix
	CONF_DEST_PATH     = "destpath"                 // destination path of your database dump files directory
	CONF_RETENTION_DAY = "retentionday"             // retention document policy day in your local storage and S3 Bucket

	CONF_REGION  = "region"  // region of your S3 Bucket
	CONF_BUCKET  = "bucket"  // bucket name of your S3 Bucket
	CONF_FOLDER  = "folder"  // folder name inside your S3 Bucket
	CONF_TIMEOUT = "timeout" // timeout when sending the data to S3 Bucket
)
