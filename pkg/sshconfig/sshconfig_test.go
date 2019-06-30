package sshconfig

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestSsh_Get(t *testing.T) {
	t.Run("default zero value", func(t *testing.T) {
		c := DefaultSsh.Get("host", "user", 22)
		assert.Equal(t, "host", c.Get("HostName"))
	})

	t.Run("custom config", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "sshconfig_tests")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		confPath := path.Join(dir, "config")
		err = ioutil.WriteFile(confPath, []byte(`
Host myhost
  ProxyCommand mycommand %h %p
`), 0600)
		assert.NoError(t, err)

		s := &Ssh{ConfigPath: confPath}
		c := s.Get("host", "user", 22)
		assert.Equal(t, "", c.Get("ProxyCommand"))

		c = s.Get("myhost", "user", 22)
		assert.Equal(t, "mycommand %h %p", c.Get("ProxyCommand"))
	})
}
