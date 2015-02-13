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
	BucketName string
}

func (p *FileUploadTask) Run(input interface{}) (interface{}, error) {
	return doUpload(p.BucketName, input.(string))
}

func doUpload(bucketName, localFilePath string) (*FileUploadResult, error) {
	//we get keys everytime because they can expire... this could be improved
	keys, err := s3gof3r.InstanceKeys() // get S3 keys
	if err != nil {
		return nil, err
	}

	s3 := s3gof3r.New(s3Root, keys)

	b := s3.Bucket(bucketName)

	file, err := os.Open(localFilePath)
	if err != nil {
		return nil, me.Err(err, "open failed")
	}
	defer file.Close()

	filename := path.Base(file.Name())
	w, err := b.PutWriter(filename, nil, nil)
	if err != nil {
		return nil, me.Err(err, "bucket writer fail")
	}

	//start := time.Now()
	size, err := io.Copy(w, file)
	if err != nil { // Copy into S3
		return nil, me.Err(err, "copy to s3 fail")
	} else {
		//u.Track.Size("s3.upload", size)
	}
	//u.Track.Duration("s3.upload", start)

	if err = w.Close(); err != nil {
		return nil, me.Err(err, "close writer fail")
	}

	os.Remove(localFilePath)
	//	os.Remove(thumbnailPath)
	//	u.Log.Debug(fmt.Sprintf("uploaded file to s3 file : %s", thumbnailPath))
	return &FileUploadResult{
		Url:      "http://" + path.Join(s3Root, bucketName, filename),
		FileSize: size,
	}, nil

}
