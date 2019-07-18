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
