package version

import (
	"fmt"
	"github.com/clarechu/infra-pulse/src/version"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"os"
)

const banner = `
   ____  __  ______  ____
  / __ \/ / / / __ \/ __ \
 / / / / /_/ / / / / / / /
/ /_/ / __  / /_/ / /_/ /
\____/_/ /_/\____/_____/

CMDB 配置管理平台
连接配置，赋能运维
`

const (
	ApiVersion = "/apis/v1alpha1"
)

func VersionCommand(args []string) *cobra.Command {
	versionCommand := &cobra.Command{
		Use:               "version",
		Short:             "CMDB 配置管理平台 ",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			info, err := version.NewBuildInfoFromOldString("")
			if err != nil {
				klog.Errorf("get version error:%s", err)
				os.Exit(-1)
			}
			fmt.Print(banner)
			fmt.Println("Version: \t" + info.Version)
			fmt.Println("GitRevision: \t" + info.GitRevision)
			fmt.Println("GolangVersion: \t" + info.GolangVersion)
			fmt.Println("BuildStatus: \t" + info.BuildStatus)
			fmt.Println("GitTag: \t" + info.GitTag)
			fmt.Println("Platform: \t" + info.Platform)
			fmt.Println("BuildDate: \t" + info.BuildDate)
		},
	}
	return versionCommand
}
