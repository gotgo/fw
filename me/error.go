package me

import (
	"math/rand"

	"github.com/krave-n/deeperror"
)

const stackFrames = 2

func Err(err error, msg string) error {
	return deeperror.NewS(rand.Int63(), msg, err, stackFrames)
}

func NewErr(msg string) error {
	return deeperror.NewS(rand.Int63(), msg, nil, stackFrames)
}
