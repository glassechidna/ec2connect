package sshconfig

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Config map[string][]string

func (s Config) Get(key string) string {
	key = strings.ToLower(key)
	vals := s[key]
	if len(vals) > 0 {
		return vals[0]
	}

	return ""
}

type Ssh struct {
	Executable string
	ConfigPath string
}

var DefaultSsh *Ssh

func (s *Ssh) Get(host, user string, port int) Config {
	args := []string{"ssh"}

	if s != nil {
		if s.Executable != "" {
			args[0] = s.Executable
		}

		if s.ConfigPath != "" {
			args = append(args, "-F", s.ConfigPath)
		}
	}

	args = append(args, "-T", "-G", "-p", strconv.Itoa(port), fmt.Sprintf("%s@%s", user, host))

	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	conf := Config{}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), " ", 2)
		key := parts[0]
		val := ""
		if len(parts) == 2 {
			val = parts[1]
		}


		conf[key] = append(conf[key], val)
	}

	return conf
}
