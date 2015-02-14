package multi_test

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"

	. "github.com/gotgo/fw/multi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var TestDownloadUrls []string

func init() {
	TestDownloadUrls = []string{
		"http://golang.org/pkg/builtin/",
		"http://golang.org/pkg/bytes/",
		"http://golang.org/pkg/compress/",
		"http://golang.org/pkg/compress/bzip2/",
		"http://golang.org/pkg/compress/flate/",
		"http://golang.org/pkg/compress/gzip/",
		"http://golang.org/pkg/compress/lzw/",
		"http://golang.org/pkg/compress/zlib/",
		"http://golang.org/pkg/container/",
		"http://golang.org/pkg/container/heap/",
		"http://golang.org/pkg/container/list/",
		"http://golang.org/pkg/container/ring/",
		"http://golang.org/pkg/crypto/",
		"http://golang.org/pkg/crypto/aes/",
		"http://golang.org/pkg/crypto/cipher/",
		"http://golang.org/pkg/crypto/des/",
		"http://golang.org/pkg/crypto/dsa/",
	}

}

var _ = Describe("FileDownload", func() {

	var (
		list       []string
		tempFolder string
		iteration  int
	)

	BeforeEach(func() {
		iteration++
		list = TestDownloadUrls
		tempFolder, _ = ioutil.TempDir("/tmp", "fileDownload")
	})

	AfterEach(func() {
		os.RemoveAll(tempFolder)
	})

	It("should download with single concurrency", func() {
		dl := &FileDownloadTask{
			Folder:  tempFolder,
			Timeout: time.Second * 10,
		}
		for i, url := range list {
			dl.Run(&FileDownloadInput{Url: url, Filename: strconv.Itoa(i)})
		}
		files, _ := ioutil.ReadDir(tempFolder)
		Expect(len(files)).To(Equal(len(list)))
	})
})
