package multi_test

import (
	"github.com/gotgo/fw/logging"
	. "github.com/gotgo/fw/multi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImageProcessor", func() {

	var (
		target *ImageProcessor
		ctx    *DataContext
	)

	BeforeEach(func() {
		target = &ImageProcessor{
			LocalPath: "/tmp/test",
			Uploader:  &NoOpUploader{},
			Log:       &logging.ConsoleLogger{},
			MaxHeight: 256,
			MaxWidth:  256,
		}
		ctx = NewDataContext()
		target.Startup()
	})

	AfterEach(func() {
		target.Shutdown()
	})

	It("single image should work", func() {
		url := "https://farm8.staticflickr.com/7327/16322277849_36e42322de_c.jpg"
		filename := "1"
		ctx.Set("data:a", "testdata")

		target.Injest(url, filename, ctx)
		result := <-target.Completed()
		Expect(result.Error).To(BeNil())
		//TODO: more tests
	})

})
