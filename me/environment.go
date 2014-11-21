package me

import (
	"fmt"
	"net"
	"os"
)

var App *Program

func init() {
	addrs, _ := net.InterfaceAddrs()
	exeName := "Not Set"
	if len(os.Args) > 0 {
		exeName = os.Args[0]
	}

	App = &Program{
		values: &ProgramValues{
			IPAddresses: addrs,
		},
		conf: &AppConf{
			Environment: "env-not-set",
			Name:        exeName,
		},
	}
}

func SetConf(c *AppConf) {
	App.conf = c
}

type AppConf struct {
	// Environment - stage, prod, dev
	Environment string `json:"environment"`
	// Name - MyAppName
	Name string `json:"name"`
}

type ProgramValues struct {
	IPAddresses []net.Addr
}

type Program struct {
	conf   *AppConf
	values *ProgramValues
}

func (p *Program) Environment() string {
	return p.conf.Environment
}

type Application interface {
	Environment() string
}

func Concat(key, seperator, value string) string {
	return fmt.Sprintf("%s%s%s", key, seperator, value)
}
