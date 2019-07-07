package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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
			err := setup("~/.ssh/config", "~/.ssh/ec2connect")
			if err != nil {
				panic(err)
			}
		},
	}

	RootCmd.AddCommand(cmd)
}

func setup(configPath, ec2connDir string) error {
	ec2connDir, err := homedir.Expand(ec2connDir)
	if err != nil {
		return err
	}

	configPath, err = homedir.Expand(configPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(ec2connDir, 0700)
	if err != nil {
		return errors.Wrapf(err, "creating directory at %s", ec2connDir)
	}

	privKey, pubKey, err := generateSshKeypair()
	if err != nil {
		return errors.Wrap(err, "generating new ssh key pair")
	}

	privKeyPath := path.Join(ec2connDir, "id_rsa")
	err = ioutil.WriteFile(privKeyPath, privKey, 0600)
	if err != nil {
		return errors.Wrap(err, "writing new ssh priv key to disk")
	}

	err = ioutil.WriteFile(path.Join(ec2connDir, "id_rsa.pub"), pubKey, 0644)
	if err != nil {
		return errors.Wrap(err, "writing new ssh pub key to disk")
	}

	snippet, err := sshConfigSnippet(privKeyPath)
	if err != nil {
		return err
	}

	myConfPath := path.Join(ec2connDir, "ssh_config")
	err = ioutil.WriteFile(myConfPath, snippet, 0644)
	if err != nil {
		return err
	}

	sshConfDir := path.Dir(configPath)
	relEc2ConfPath, err := filepath.Rel(sshConfDir, myConfPath)
	if err != nil {
		return err
	}

	err = idempotentInsert(configPath, fmt.Sprintf("Include %s\n\n", relEc2ConfPath))
	return errors.Wrapf(err, "appending config to %s", configPath)
}

func sshConfigSnippet(privKeyPath string) ([]byte, error) {
	cmdPath, err := exec.LookPath("ec2connect")
	if err != nil {
		return nil, errors.Wrap(err, "You have to first install ec2connect somewhere on your PATH")
	}

	tmpl, err := template.New("").Parse(`
Match exec "{{ .CommandPath }} match --host %n --user %r"
  IdentityFile {{ .KeyPath }}
  ProxyCommand {{ .CommandPath }} connect --instance-id %h --user %r --port %p
`)
	if err != nil {
		return nil, errors.Wrap(err, "parsing ssh config template")
	}

	b := &bytes.Buffer{}
	err = tmpl.Execute(b, map[string]string{
		"KeyPath":     privKeyPath,
		"CommandPath": cmdPath,
	})
	if err != nil {
		return nil, errors.Wrap(err, "rendering ssh config template")
	}

	return b.Bytes(), nil
}

func idempotentInsert(path, content string) error {
	existing, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "reading existing file")
	}

	if strings.Contains(string(existing), content) {
		return nil
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrap(err, "opening file for appending")
	}
	defer f.Close()

	_, err = f.WriteString(content)
	_, err = f.Write(existing)
	return err
}

func generateSshKeypair() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, errors.Wrap(err, "generating new ssh private key")
	}

	buf := bytes.Buffer{}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	err = pem.Encode(&buf, privateKeyPEM)
	if err != nil {
		return nil, nil, errors.Wrap(err, "pem-encoding new ssh private key")
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "creating signer from ssh priv key")
	}

	public := ssh.MarshalAuthorizedKey(signer.PublicKey())
	return buf.Bytes(), public, nil
}
