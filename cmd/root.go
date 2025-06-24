package cmd

import (
	"github.com/clarechu/infra-pulse/cmd/version"
	"github.com/spf13/cobra"
)

func GetRootCmd(args []string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:               "cmdb",
		Short:             "cmdb  interface.",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Long:              `The new generation of CMDB`,
	}
	rootCmd.AddCommand(ServerCommand(args))
	rootCmd.AddCommand(version.VersionCommand(args))

	return rootCmd
}
