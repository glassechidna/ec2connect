package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"regexp"
)

func init() {
	cmd := &cobra.Command{
		Use:   "match",
		Short: "Internal functionality",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.PersistentFlags().GetString("host")
			user, _ := cmd.PersistentFlags().GetString("user")
			matched := match(host, user)
			if !matched {
				os.Exit(1)
			}
		},
	}

	cmd.PersistentFlags().String("host", "", "")
	cmd.PersistentFlags().String("user", "", "")
	RootCmd.AddCommand(cmd)
}

func match(host, user string) bool {
	re := regexp.MustCompile(`^i-[a-f0-9]+`)
	doesMatch := re.MatchString(host)

	//if doesMatch {
	//	dir, _ := homedir.Expand("~/.ssh/ec2connect")
	//	filePath := path.Join(dir, host)
	//	ioutil.WriteFile(filePath, []byte(user), 0644)
	//}

	return doesMatch
}
