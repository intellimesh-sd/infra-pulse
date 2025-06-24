package version

import (
	"fmt"
	"strings"
)

// The following fields are populated at build time using -ldflags -X.
// Note that DATE is omitted for reproducible builds
var (
	buildVersion     = "unknown"
	buildGitRevision = "unknown"
	buildStatus      = "unknown"
	buildTag         = "unknown"
	golangVersion    = "unknown"
	buildHub         = "unknown"
	buildDate        = "unknown"
	platform         = "unknown"
)

// BuildInfo describes version information about the binary build.
type BuildInfo struct {
	Version       string `json:"version"`
	GitRevision   string `json:"revision"`
	GolangVersion string `json:"golangVersion"`
	BuildStatus   string `json:"status"`
	BuildDate     string `json:"buildDate"`
	GitTag        string `json:"tag"`
	Platform      string `json:"platform"`
}

// ServerInfo contains the version for a single control plane component
type ServerInfo struct {
	Component string
	Info      BuildInfo
}

// MeshInfo contains the versions for all Istio control plane components
type MeshInfo []ServerInfo

// ProxyInfo contains the version for a single data plane component
type ProxyInfo struct {
	ID           string
	IstioVersion string
}

// DockerBuildInfo contains and exposes Hub: buildHub and Tag: buildVersion
type DockerBuildInfo struct {
	Hub string
	Tag string
}

// NewBuildInfoFromOldString creates a BuildInfo struct based on the output
// of previous olm components '-- version' output
func NewBuildInfoFromOldString(oldOutput string) (BuildInfo, error) {
	res := Info

	lines := strings.Split(oldOutput, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.SplitN(line, ":", 2)
		if fields != nil {
			if len(fields) != 2 {
				return BuildInfo{}, fmt.Errorf("invalid BuildInfo input, field '%s' is not valid", fields[0])
			}
			value := strings.TrimSpace(fields[1])
			switch fields[0] {
			case "Version":
				res.Version = value
			case "GitRevision":
				res.GitRevision = value
			case "GolangVersion":
				res.GolangVersion = value
			case "BuildStatus":
				res.BuildStatus = value
			case "GitTag":
				res.GitTag = value
			case "BuildDate":
				res.BuildDate = value
			case "Platform":
				res.Platform = value
			default:
				// Skip unknown fields, as older versions may report other fields
				continue
			}
		}
	}

	return res, nil
}

var (
	// Info exports the build version information.
	Info       BuildInfo
	DockerInfo DockerBuildInfo
)

// String produces a single-line version info
//
// This looks like:
//
// ```
// u<version>-<git revision>-<build status>
// ```
func (b BuildInfo) String() string {
	return fmt.Sprintf(`Version:%v GIT_REVISTION:%v BUILD_STATUS:%v`,
		b.Version,
		b.GitRevision,
		b.BuildStatus)
}

// LongForm returns a dump of the Info struct
// This looks like:
func (b BuildInfo) LongForm() string {
	return fmt.Sprintf("%#v", b)
}

func init() {
	Info = BuildInfo{
		Version:       buildVersion,
		GitRevision:   buildGitRevision,
		GolangVersion: golangVersion,
		BuildStatus:   buildStatus,
		GitTag:        buildTag,
		BuildDate:     buildDate,
		Platform:      platform,
	}

	DockerInfo = DockerBuildInfo{
		Hub: buildHub,
		Tag: buildVersion,
	}
}
