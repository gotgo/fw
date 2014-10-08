package io

import (
	"errors"
	"net"
	"time"
)

var StoppedError = errors.New("Listener stopped")

type GracefulListener struct {
	*net.TCPListener
	shutdown chan int
}

func MakeGraceful(l net.Listener) (*GracefulListener, error) {
	if listener, ok := l.(*net.TCPListener); !ok {
		return nil, errors.New("Cannot wrap listener")
	} else {
		retval := &GracefulListener{
			TCPListener: listener,
			shutdown:    make(chan int),
		}
		return retval, nil
	}
}

func (gl *GracefulListener) Accept() (net.Conn, error) {

	for {
		//Wait up to one second for a new connection
		gl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := gl.TCPListener.Accept()

		//Check for the channel being closed
		select {
		case <-gl.shutdown:
			return nil, StoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (gl *GracefulListener) Shutdown() {
	close(gl.shutdown)
}
