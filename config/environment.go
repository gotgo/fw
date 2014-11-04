package config

import (
	"fmt"
	"net"
	"os"
)

var Prog *Program

func init() {
	hostname, _ := os.Hostname()
	addrs, _ := net.InterfaceAddrs()
	envVars := os.Environ()
	exeName := "Not Set"
	if len(os.Args) > 0 {
		exeName = os.Args[0]
	}

	Prog = &Program{
		Environment:  "stage",
		HostName:     hostname,
		IPAddresses:  addrs,
		Name:         exeName,
		EnvVariables: envVars,
	}
}

type Program struct {
	Environment  string
	HostName     string
	IPAddresses  []net.Addr
	Name         string
	EnvVariables []string
}

func (p *Program) AppendEnvironment(key string, seperator string) string {
	return fmt.Sprintf("%s%s%s", key, seperator, p.Environment)
}
