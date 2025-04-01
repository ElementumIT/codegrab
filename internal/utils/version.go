package utils

var (
	Version   = "dev"
	CommitSHA = "unknown"
)

func VersionInfo() string {
	return "grab " + Version + " (" + CommitSHA + ")"
}
