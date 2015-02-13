package multi_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	. "github.com/gotgo/fw/multi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FileDownload", func() {

	var (
		list       []string
		tempFolder string
		iteration  int
	)

	BeforeEach(func() {
		iteration++

		tempFolder, _ = ioutil.TempDir("/tmp", "fileDownload")

		list = []string{
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
	})

	AfterEach(func() {
		os.RemoveAll(tempFolder)
	})

	It("should download with single concurrency", func() {
		dl := &FileDownloads{
			Folder:       tempFolder,
			Concurrency:  1,
			MaxQueuedIn:  1,
			MaxQueuedOut: 1,
		}

		for i, url := range list {
			fmt.Printf("\nadd %d\n", i)
			dl.Add(&FileDownloadsInput{Url: url, Filename: strconv.Itoa(i)})
			<-dl.Completed()
		}

		files, _ := ioutil.ReadDir(tempFolder)
		Expect(len(files)).To(Equal(len(list)))
	})
	It("should download with multiple concurrency", func() {
		fmt.Println("hello 2")
		dl := &FileDownloads{
			Folder:       tempFolder,
			Concurrency:  4,
			MaxQueuedIn:  8,
			MaxQueuedOut: len(list),
		}

		for i, url := range list {
			fmt.Printf("\nadd %d\n", i)
			dl.Add(&FileDownloadsInput{Url: url, Filename: strconv.Itoa(i)})
		}

		dl.Shutdown()

		i := 0
		for c := range dl.Completed() {
			i++
			fmt.Printf("\ncomplete %d %s\n", i, c.Input.Url)
		}

		files, _ := ioutil.ReadDir(tempFolder)
		Expect(len(files)).To(Equal(len(list)))
	})

})
