package sshconfig

import (
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"strings"
)

func (c Config) IdentityAgent() agent.Agent {
	conf := c.Get("IdentityAgent")
	if conf == "none" {
		return nil
	}

	if conf == "" || conf == "SSH_AUTH_SOCK" {
		return agentAtPath(os.Getenv("SSH_AUTH_SOCK"))
	}

	conf = c.Expand(conf)

	if strings.HasPrefix(conf, "$") {
		return agentAtPath(os.Getenv(conf[1:]))
	}

	return agentAtPath(conf)
}

func agentAtPath(path string) agent.Agent {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil
	}

	return agent.NewClient(conn)
}
