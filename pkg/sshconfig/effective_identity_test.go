package sshconfig

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestConfig_EffectivePublicKey(t *testing.T) {
	t.Run("only choose pubkey when corresponding priv key exists", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "sshconfig_tests")
		assert.NoError(t, err)
		defer os.RemoveAll(dir)

		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		assert.NoError(t, err)

		buf := bytes.Buffer{}
		privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
		err = pem.Encode(&buf, privateKeyPEM)
		assert.NoError(t, err)

		signer, err := ssh.NewSignerFromKey(privateKey)
		assert.NoError(t, err)
		pubKey := ssh.MarshalAuthorizedKey(signer.PublicKey())

		err = ioutil.WriteFile(path.Join(dir, "id_rsa.pub"), pubKey, 0644)
		assert.NoError(t, err)

		confPath := path.Join(dir, "ssh_config")
		err = ioutil.WriteFile(confPath, []byte(fmt.Sprintf(`
IdentityFile %s
IdentitiesOnly yes
`, path.Join(dir, "id_rsa"))), 0644)
		assert.NoError(t, err)

		s := &Ssh{ConfigPath: confPath}
		c := s.Get("any", "user", 22)
		key := c.EffectivePublicKey()
		assert.Nil(t, key)
	})
}