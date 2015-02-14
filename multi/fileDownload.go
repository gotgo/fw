package multi

import (
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/gotgo/fw/me"
)

func download(url, filename, folder string, timeout time.Duration) (*FileDownloadOutput, error) {
	fp := path.Join(folder, filename)
	output := &FileDownloadOutput{
		Path: fp,
	}

	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, me.Err(err, "failed to get URL ")
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, me.Err(err, "failed to read all bytes")
	}

	resp.Body.Close()

	if err := ioutil.WriteFile(fp, bts, 0666); err != nil {
		return nil, me.Err(err, "failed to write file to "+fp)
	}

	//output.ContentType = resp.ContentType
	//d.Track.Duration("download", started)
	//d.Track.Size("download", resp.ContentLength)

	return output, nil
}

const defaultTimeout = time.Second * 60

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
