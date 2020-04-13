package cli

const (
	// Name is the current application name
	Name = "wg-registry"

	// Version is the current application version
	Version = "0.1.0"
)

var (
	// GitCommit contains the Git commit SHA for the binary
	GitCommit string

	// BuildTime contains the binary build time
	BuildTime string

	// GoVersion contains the Go runtime version
	GoVersion string
)
