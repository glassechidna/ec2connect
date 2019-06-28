package cmd

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

func init() {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "First-time setup of ec2connect on your machine",
		Long: `
'setup' will configure your ~/.ssh/config to use the 'ec2connect' helper tool
whenever you ssh into an EC2 server. This setup only needs to be run once on
your machine.
`,
		Run: func(cmd *cobra.Command, args []string) {
			configPath, _ := cmd.PersistentFlags().GetString("config-path")
			keyPath, _ := cmd.PersistentFlags().GetString("key-path")
			err := setup(configPath, keyPath)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.PersistentFlags().String("config-path", "~/.ssh/config", "")
	cmd.PersistentFlags().String("key-path", "~/.ssh/id_rsa", "")
	RootCmd.AddCommand(cmd)
}

func setup(configPath, keyPath string) error {
	tmpl, err := template.New("").Parse(`
Match exec "ec2connect match --host %n --user %r"
  IdentityFile {{ .KeyPath }}
  ProxyCommand ec2connect connect --instance-id %h --user %r --ssh-key {{ .KeyPath }}
`)
	if err != nil {
		return err
	}

	keyPath, err = homedir.Expand(keyPath)
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	err = tmpl.Execute(b, map[string]string{"KeyPath": keyPath})
	if err != nil {
		return err
	}

	configPath, err = homedir.Expand(configPath)
	if err != nil {
		return err
	}

	confDir := path.Dir(configPath)
	ec2connDir := path.Join(confDir, "ec2connect")
	err = os.MkdirAll(ec2connDir, 0644)
	if err != nil {
		return err
	}

	myConfPath := path.Join(confDir, "ec2connect_config")
	err = ioutil.WriteFile(myConfPath, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf(`

Include %s
`, myConfPath))
	return err
}
