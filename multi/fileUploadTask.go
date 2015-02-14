package multi

import (
	"io"
	"os"
	"path"

	"github.com/gotgo/fw/me"
	"github.com/rlmcpherson/s3gof3r"
)

const s3Root = "s3-us-west-2.amazonaws.com"

type FileUploadResult struct {
	Url      string
	FileSize int64
}

type FileUploadTask struct {
	Uploader Uploader
}

func (u *FileUploadTask) Run(input interface{}) (interface{}, error) {
	return doUpload(u.Uploader, input.(string))
}

func (u *FileUploadTask) Name() string {
	return "upload"
}

func doUpload(uploader Uploader, localFilePath string) (*FileUploadResult, error) {

	file, err := os.Open(localFilePath)
	if err != nil {
		return nil, me.Err(err, "open failed")
	}
	defer file.Close()
	filename := path.Base(file.Name())

	w, err := uploader.Writer(filename)
	if err != nil {
		return nil, me.Err(err, "failed to get writer")
	}

	size, err := io.Copy(w, file)
	if err != nil { // Copy into S3
		return nil, me.Err(err, "copy to s3 fail")
	}

	if err = w.Close(); err != nil {
		return nil, me.Err(err, "close writer fail")
	}

	os.Remove(localFilePath)
	//	os.Remove(thumbnailPath)
	//	u.Log.Debug(fmt.Sprintf("uploaded file to s3 file : %s", thumbnailPath))
	return &FileUploadResult{
		Url:      uploader.DestinationUrl(filename),
		FileSize: size,
	}, nil

}

type NoOpUploader struct {
}

func (w *NoOpUploader) Writer(filename string) (io.WriteCloser, error) {
	return &NoOpWriteCloser{}, nil
}

func (w *NoOpUploader) DestinationUrl(filename string) string {
	return filename
}

type NoOpWriteCloser struct {
}

func (n *NoOpWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (n *NoOpWriteCloser) Close() error {
	return nil
}

type S3Writer struct {
	BucketName string
	Root       string
}

func (w *S3Writer) Writer(filename string) (io.WriteCloser, error) {
	//we get keys everytime because they can expire... this could be improved
	keys, err := s3gof3r.InstanceKeys() // get S3 keys
	if err != nil {
		return nil, err
	}

	s3 := s3gof3r.New(w.Root, keys)
	b := s3.Bucket(w.BucketName)

	return b.PutWriter(filename, nil, nil)
}

func (w *S3Writer) DestinationUrl(filename string) string {
	return "http://" + path.Join(w.Root, w.BucketName, filename)
}

type Uploader interface {
	Writer(string) (io.WriteCloser, error)
	DestinationUrl(string) string
}
