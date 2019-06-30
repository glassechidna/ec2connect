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

		privBytes, err := ioutil.ReadFile(identityPath)
		if err != nil {
			continue
		}

		pubBytes, err := ioutil.ReadFile(identityPath + ".pub")
		if err == nil {
			return pubBytes
		}

		signer, err := ssh.ParsePrivateKey(privBytes)
		if err != nil {
			continue
		}

		return ssh.MarshalAuthorizedKey(signer.PublicKey())
	}

	if c.Get("IdentitiesOnly") == "yes" {
		return nil
	}

	agent := c.IdentityAgent()
	if agent == nil {
		return nil
	}

	keys, err := agent.List()
	if err != nil || len(keys) == 0 {
		return nil
	}

	return ssh.MarshalAuthorizedKey(keys[0])
}
