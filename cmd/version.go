package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// these are set by goreleaser
var version, commit, date string

func init() {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Information about this build of ec2connect",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`
Version: %s
Commit: %s
Date: %s
`, version, commit, date)
		},
	}

	RootCmd.AddCommand(cmd)
}
