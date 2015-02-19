package multi_test

import (
	"bytes"
	"io/ioutil"
	"os"

	. "github.com/gotgo/fw/multi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AcquireImage", func() {

	It("should resize", func() {
		file, err := os.Open("./img/cindy1.jpg")
		if err != nil {
			Expect(err.Error()).To(Equal(""))
		}
		defer file.Close()
		bts, err := ioutil.ReadAll(file)
		if err != nil {
			Expect(err.Error()).To(Equal(""))
		}
		a := &AcquireImage{}
		_, _, _, err = a.Resize(bytes.NewReader(bts), 256, 256)
		if err != nil {
			Expect(err.Error()).To(Equal(""))
		}
	})
})
