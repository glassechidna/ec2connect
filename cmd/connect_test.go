package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestAuthorize(t *testing.T) {
	instanceId := os.Getenv("TEST_INSTANCE_ID")

	if testing.Short() || instanceId == "" {
		t.SkipNow()
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	buf := bytes.Buffer{}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	err = pem.Encode(&buf, privateKeyPEM)
	assert.NoError(t, err)

	err = ioutil.WriteFile("/tmp/sshkey", buf.Bytes(), 0600)
	assert.NoError(t, err)

	info, err := authorize(instanceId, "ap-southeast-2", "ec2-user", "/tmp/sshkey")
	assert.NoError(t, err)

	assert.True(t, len(info.Address) > 0)
}
