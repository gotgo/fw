package io_test

import (
	"net"
	"net/http"
	"sync"
	"time"

	. "github.com/gotgo/fw/io"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

var _ = Describe("GracefulListener", func() {
	//fn' mac launches a warning window everytime this is run
	It("should stop", func() {
		originalListener, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}

		listener, err := MakeGraceful(originalListener)
		if err != nil {
			panic(err)
		}

		smux := http.NewServeMux()
		smux.HandleFunc("/", testHttp)
		server := http.Server{
			Handler: smux,
		}

		var wg sync.WaitGroup
		go func() {
			wg.Add(1)
			defer wg.Done()
			server.Serve(listener)
		}()

		start := time.Now()
		listener.Shutdown()
		wg.Wait()
		duration := time.Since(start)
		twoSec := 3 * time.Second
		Expect(duration).Should(BeNumerically("<", twoSec))
	})

})
