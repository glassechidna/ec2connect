package sshconfig

import (
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func (c Config) EffectivePublicKey() []byte {
	for _, identityPath := range c["identityfile"] {
		identityPath, err := homedir.Expand(identityPath)
		if err != nil {
			continue
		}

		pubBytes, err := ioutil.ReadFile(identityPath + ".pub")
		if err == nil {
			return pubBytes
		}

		privBytes, err := ioutil.ReadFile(identityPath)
		if err != nil {
			continue
		}

		signer, err := ssh.ParsePrivateKey(privBytes)
		if err != nil {
			continue
		}

		return ssh.MarshalAuthorizedKey(signer.PublicKey())
	}

	// TODO support ssh agents
	return nil
}
