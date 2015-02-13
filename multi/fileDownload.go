package multi

import (
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gotgo/fw/logging"
	"github.com/gotgo/fw/me"
	"github.com/gotgo/fw/stats"
)

const downloadTimeOut = time.Second * 60

func NewFileDownloads(folder string, concurrency int, inMax, outMax int) *FileDownloads {
	fd := &FileDownloads{
		Folder:      folder,
		Concurrency: concurrency,
	}
	return fd
}

func (d *FileDownloads) setup() {
	if d.Folder == "" {
		d.Folder = os.TempDir()
	}
	if d.Concurrency == 0 {
		d.Concurrency = 1
	}
	if d.MaxQueuedIn == 0 {
		d.MaxQueuedIn = 1
	}
	if d.MaxQueuedOut == 0 {
		d.MaxQueuedOut = 1
	}
	d.input = make(chan *FileDownloadsInput, d.MaxQueuedIn)
	d.output = make(chan *FileDownloadsResult, d.MaxQueuedOut)
	d.shutdown = make(chan struct{})
	d.outstanding = &sync.WaitGroup{}
	d.once = &sync.Once{}

	d.Track = stats.NewBasicMeter("image.downloader", me.App.Environment())
}

func (d *FileDownloads) Start() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.running {
		panic("already running")
	}
	d.running = true
	//TODO: prevent start from being called twice
	d.setup()
	concurrency := d.Concurrency

	for i := 0; i < concurrency; i++ {
		go d.run()
	}
}

type FileDownloads struct {
	Folder string

	mutex        sync.Mutex
	Track        stats.BasicMeter
	Log          logging.Logger
	Concurrency  int
	MaxQueuedIn  int
	MaxQueuedOut int
	running      bool
	input        chan *FileDownloadsInput
	output       chan *FileDownloadsResult
	outstanding  *sync.WaitGroup
	once         *sync.Once
	done         chan struct{}
	shutdown     chan struct{}
}

type FileDownloadsInput struct {
	Url      string
	Filename string
}

type FileDownloadsResult struct {
	Error  error
	Input  *FileDownloadsInput
	Output *FileDownloadsOutput
}

type FileDownloadsOutput struct {
	Path        string
	Size        int64
	ContentType string
}

// Add - will block when the number of items queued reaches MaxQueuedInput
func (d *FileDownloads) Add(todo *FileDownloadsInput) {
	d.input <- todo
}

func (d *FileDownloads) Completed() <-chan *FileDownloadsResult {
	return d.output
}

// Shutdown - begins the shutdown operation. Reading on the returned channel will block until the shutdown is complete
func (d *FileDownloads) Shutdown() chan struct{} {
	close(d.input) //no more input
	return d.done
}

func (d *FileDownloads) run() {
	for in := range d.input {
		d.safeDownload(in)
	}
	d.outstanding.Wait()

	d.once.Do(func() {
		close(d.output)
		close(d.done)
	})
}

func (d *FileDownloads) safeDownload(in *FileDownloadsInput) {
	folder := d.Folder
	d.outstanding.Add(1)
	defer func() {
		d.outstanding.Done()
		if r := recover(); r != nil {
			me.LogRecoveredPanic(d.Log, "download failed", r, &logging.KV{"from", d})
		}
	}()

	out, err := download(in.Url, in.Filename, folder)
	d.output <- &FileDownloadsResult{
		Error:  err,
		Input:  in,
		Output: out,
	}
}

func download(url, filename, folder string) (*FileDownloadsOutput, error) {
	output := &FileDownloadsOutput{}

	//create file first, so we know we're able to save to disk
	fp := path.Join(folder, filename)
	file, err := os.Create(fp)
	if err != nil {
		return nil, me.Err(err, "create file")
	}
	defer file.Close()

	//d.Log.Debug("downloading " + fp)

	//download
	//started := time.Now()
	client := http.Client{
		Timeout: downloadTimeOut,
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

	output.Path = fp
	return output, nil
}

type DownloadTask struct {
	Folder string
}

func (d *DownloadTask) Run(input interface{}) (interface{}, error) {
	in, ok := input.(*FileDownloadsInput)
	if !ok {
		panic("unexpected type")
	}
	return download(in.Url, in.Filename, d.Folder)
}
