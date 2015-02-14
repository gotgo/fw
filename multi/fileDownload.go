package multi

import (
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gotgo/fw/me"
)

func download(url, filename, folder string, timeout time.Duration) (*FileDownloadOutput, error) {
	output := &FileDownloadOutput{}

	//create file first, so we know we're able to save to disk
	fp := path.Join(folder, filename)
	file, err := os.Create(fp)
	if err != nil {
		return nil, me.Err(err, "create file")
	}
	defer file.Close()

	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, me.Err(err, "download")
	}
	defer resp.Body.Close()
	//output.ContentType = resp.ContentType
	//d.Track.Duration("download", started)
	//d.Track.Size("download", resp.ContentLength)

	//save to disk
	if size, err := io.Copy(file, resp.Body); err != nil {
		return nil, me.Err(err, "failed to save downloaded file")
	} else if size > 0 {
		output.Size = size
		//	d.Track.Size("saved", size)
	}
	file.Sync()
	file.Close()
	output.Path = fp
	return output, nil
}

const defaultTimeout = time.Second * 30

type FileDownloadTask struct {
	Folder  string
	Timeout time.Duration
}

func (d *FileDownloadTask) Run(input interface{}) (interface{}, error) {
	in, ok := input.(*FileDownloadInput)
	if !ok {
		panic("unexpected type")
	}
	timeout := d.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return download(in.Url, in.Filename, d.Folder, timeout)
}

func (d *FileDownloadTask) Name() string {
	return "fileDownload"
}

type FileDownloadInput struct {
	Url      string
	Filename string
}

type FileDownloadOutput struct {
	Path        string
	Size        int64
	ContentType string
}
