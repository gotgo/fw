package multi

import (
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/gotgo/fw/logging"
)

type ImageProcessor struct {

	// Uploader - required
	Uploader Uploader

	// LocalPath - defaults to os.TempDir()/downloads
	LocalPath string

	Log logging.Logger

	MaxHeight int
	MaxWidth  int

	downloader *TaskRun
	phasher    *TaskRun
	resizer    *TaskRun
	uploader   *TaskRun

	outstanding sync.WaitGroup
	complete    chan *ImageProcessorOutput
}

func (ip *ImageProcessor) setup() {
	uploader := ip.Uploader
	if uploader == nil {
		panic("uploader not set")
	}

	if ip.Log == nil {
		ip.Log = &logging.NoOpLogger{}
	}

	var tempFolder = ip.LocalPath
	if tempFolder == "" {
		tempFolder = path.Join(os.TempDir(), "downloads")
	}

	os.RemoveAll(tempFolder)
	os.MkdirAll(tempFolder, 0777)

	ip.complete = make(chan *ImageProcessorOutput, 100)

	ip.downloader = &TaskRun{
		Action:       &FileDownloadTask{Folder: tempFolder},
		Concurrency:  1,
		MaxQueuedIn:  10 * 5,
		MaxQueuedOut: 10 * 10,
	}

	ip.phasher = &TaskRun{
		Action: &PHashTask{
			Log: ip.Log,
		},
		Concurrency:  1,
		MaxQueuedIn:  2,
		MaxQueuedOut: 100,
	}

	ip.resizer = &TaskRun{
		Action:       &ResizeImageTask{MaxHeight: ip.MaxHeight, MaxWidth: ip.MaxWidth},
		Concurrency:  1,
		MaxQueuedIn:  12,
		MaxQueuedOut: 100,
	}

	ip.uploader = &TaskRun{
		Action:       &FileUploadTask{Uploader: uploader},
		Concurrency:  1,
		MaxQueuedIn:  8 * 10,
		MaxQueuedOut: 100,
	}
}

func (p *ImageProcessor) Startup() {
	p.setup()
	p.doScavenge()
}

func (p *ImageProcessor) Shutdown() {
	p.downloader.Shutdown()
	p.outstanding.Wait()
}

func (p *ImageProcessor) Injest(url, filename string, ctx *DataContext) {
	in := &FileDownloadInput{
		Url:      url,
		Filename: filename,
	}

	p.downloader.Add(in, ctx)
}

func (p *ImageProcessor) doScavenge() {
	p.downloader.Startup()
	p.phasher.Startup()
	p.resizer.Startup()
	p.uploader.Startup()

	p.outstanding.Add(1)
	// phash
	go p.phash()
	p.outstanding.Add(1)
	// resize
	go p.resize()
	p.outstanding.Add(1)
	// upload
	go p.upload()
	p.outstanding.Add(1)
	// collector
	go p.wrapUp()
}

func (p *ImageProcessor) handleError(message string, result *TaskRunOutput) {
	p.Log.Error(message, result.Error())
	p.complete <- &ImageProcessorOutput{
		Error:   result.Error(),
		Context: result.Context,
	}
}

func (p *ImageProcessor) phash() {
	for dl := range p.downloader.Completed() {
		mutex.Lock()
		mutex.Unlock()

		if dl.Error() != nil {
			p.handleError("download failed", dl)
		} else {
			result := dl.Context.Get(p.downloader.Name()).(*TaskRunResult)
			dlo := result.Output.(*FileDownloadOutput)
			p.phasher.Add(dlo.Path, dl.Context)
		}
	}
	p.phasher.Shutdown()
	p.outstanding.Done()
}

func (p *ImageProcessor) resize() {
	i := 0
	for ph := range p.phasher.Completed() {
		mutex.Lock()
		mutex.Unlock()

		i++
		result := ph.Context.Get(p.downloader.Name()).(*TaskRunResult)

		if ph.Error() != nil {
			dli := result.Input.(*FileDownloadInput)
			p.handleError("Phash failed ct:"+strconv.Itoa(i)+" "+dli.Url, ph)
		} else {
			dlo := result.Output.(*FileDownloadOutput)
			p.resizer.Add(dlo.Path, ph.Context)
		}
	}
	p.resizer.Shutdown()
	p.outstanding.Done()
}

func (p *ImageProcessor) upload() {
	for rz := range p.resizer.Completed() {
		mutex.Lock()
		mutex.Unlock()

		if rz.Error() != nil {
			p.handleError("resize failed", rz)
		} else {
			rzOut := rz.Output().(*ImageResizeOutput)
			p.uploader.Add(rzOut.FilePath, rz.Context)
		}
	}
	p.uploader.Shutdown()
	p.outstanding.Done()
}

var mutex sync.Mutex

func (p *ImageProcessor) wrapUp() {
	for result := range p.uploader.Completed() {
		mutex.Lock()
		mutex.Unlock()

		if result.Error() != nil {
			p.handleError("upload failed", result)
		} else {
			dl := result.Previous(p.downloader.Name()).Output.(*FileDownloadOutput)
			dlin := result.Previous(p.downloader.Name()).Input.(*FileDownloadInput)
			rz := result.Previous(p.resizer.Name()).Output.(*ImageResizeOutput)
			phash, _ := result.Previous(p.phasher.Name()).Output.(uint64)
			ul := result.Previous(p.uploader.Name()).Output.(*FileUploadOutput)

			if dl == nil || dlin == nil {
				panic("nil - download in or out")
			}
			if rz == nil {
				r := result.Context.Get(p.resizer.Name()).(*TaskRunResult)
				if r.Error != nil {
					p.Log.Error("slipped through error", r.Error)
				} else if r.Output != nil {
					p.Log.Warn("WTF")
				} else {
					panic("nil - resizer")
				}
			}
			if ul == nil {
				panic("nil - upload")
			}
			if result == nil {
				panic("nil - result")
			}

			r := &ImageProcessorOutput{
				DownloadSize:        dl.Size,
				DownloadContentType: dl.ContentType,
				DownloadUrl:         dlin.Url,
				PHash:               phash,
				FileSize:            rz.FileSize,
				Height:              rz.Height,
				Width:               rz.Width,
				ContentType:         rz.ContentType,
				DestinationUrl:      ul.Url,
				Error:               result.Error(),
				Context:             result.Context,
			}

			p.complete <- r
		}
	}
	close(p.complete)
	p.outstanding.Done()
}

func (p *ImageProcessor) Completed() <-chan *ImageProcessorOutput {
	return p.complete
}

type ImageProcessorOutput struct {
	DownloadSize        int64
	DownloadContentType string
	DownloadUrl         string
	PHash               uint64

	FileSize    int64
	Height      int
	Width       int
	ContentType string

	DestinationUrl string
	Error          error
	Context        *DataContext
}
