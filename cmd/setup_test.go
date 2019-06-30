package cmd

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"testing"
)

func TestSetup(t *testing.T) {
	t.Run("ssh_config doesn't exist", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "ec2connect_tests")
		assert.NoError(t, err)

		fmt.Println(dir)

		sshConfPath := path.Join(dir, "ssh_config")
		ec2ConnDir := path.Join(dir, "ec2connect")
		err = setup(sshConfPath, ec2ConnDir)
		assert.NoError(t, err)
	})
}
