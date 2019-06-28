package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Remotely install EC2 Instance Connect on EC2 instance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	RootCmd.AddCommand(cmd)
}

func install() error {
	return nil
}
