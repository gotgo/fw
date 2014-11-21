package me

import (
	"math/rand"

	"github.com/amattn/deeperror"
)

func Err(err error, msg string) error {
	//TODO: use stack frame version
	return deeperror.New(rand.Int63(), msg, err)
}
