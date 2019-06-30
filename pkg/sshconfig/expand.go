package sshconfig

import (
	"os"
	"os/user"
	"strings"
)

func (c Config) Expand(input string) string {
	r := strings.ReplaceAll
	s := input

	u, _ := user.Current()
	host, _ := os.Hostname()
	remote := c.Get("HostName")

	s = r(s, "%C", "") // TODO
	s = r(s, "%d", u.HomeDir)
	s = r(s, "%h", remote)
	s = r(s, "%i", u.Uid)
	s = r(s, "%L", host)
	s = r(s, "%l", host)
	s = r(s, "%n", remote)
	s = r(s, "%p", c.Get("Port"))
	s = r(s, "%r", c.Get("User"))
	s = r(s, "%T", "NONE")
	s = r(s, "%u", u.Username)
	s = r(s, "%%", "%")

	return s
}
