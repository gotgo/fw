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

	App = NewProgram(&ProgramValues{
		Environment: "env-not-set",
		IPAddresses: addrs,
		Name:        exeName,
	})

}

type ProgramValues struct {
	Environment string
	IPAddresses []net.Addr
	Name        string
}

type Program struct {
	values *ProgramValues
}

func NewProgram(values *ProgramValues) *Program {
	p := &Program{
		values: values,
	}
	return p
}

func (p *Program) Environment() string {
	return p.values.Environment
}

type Application interface {
	Environment() string
}

func Concat(key, seperator, value string) string {
	return fmt.Sprintf("%s%s%s", key, seperator, value)
}
