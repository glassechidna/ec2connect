package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
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

	signer, err := ssh.NewSignerFromKey(privateKey)
	assert.NoError(t, err)

	pubKey := ssh.MarshalAuthorizedKey(signer.PublicKey())
	info, err := authorize(instanceId, "ap-southeast-2", "ec2-user", string(pubKey))
	assert.NoError(t, err)

	assert.True(t, len(info.Address) > 0)
}
