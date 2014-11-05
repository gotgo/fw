package config

import (
	"fmt"
	"net"
	"os"
)

var App *Application

func init() {
	hostname, _ := os.Hostname()
	addrs, _ := net.InterfaceAddrs()
	envVars := os.Environ()
	exeName := "Not Set"
	if len(os.Args) > 0 {
		exeName = os.Args[0]
	}

	App = &Application{
		Environment:  "stage",
		HostName:     hostname,
		IPAddresses:  addrs,
		Name:         exeName,
		EnvVariables: envVars,
	}
}

type Application struct {
	Environment  string
	HostName     string
	IPAddresses  []net.Addr
	Name         string
	EnvVariables []string
}

func (a *Application) AppendEnvironment(key string, seperator string) string {
	return fmt.Sprintf("%s%s%s", key, seperator, a.Environment)
}
