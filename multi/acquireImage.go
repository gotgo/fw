package multi

import (
	"bytes"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"

	"github.com/daddye/vips"
	"github.com/disintegration/imaging"
	"github.com/gotgo/fw/logging"
	"github.com/gotgo/fw/me"
)

type AcquireImage struct {
	DownloadTimeout time.Duration
	MaxHeight       int
	MaxWidth        int
	Uploader        Uploader
	Log             logging.Logger
}

type AcquiredImage struct {
	SourceSize        int
	SourceUrl         string
	SourceContentType string

	DestUrl    string
	DestSize   int64
	DestHeight int
	DestWidth  int
	PHash      uint64
}

func (a *AcquireImage) timeout() time.Duration {
	if a.DownloadTimeout == time.Second*0 {
		return time.Second * 60
	}
	return a.DownloadTimeout
}

func (a *AcquireImage) Acquire(url, filename string) (*AcquiredImage, error) {
	//TODO: time
	bts, ctype, err := a.download(url, a.timeout())
	if err != nil {
		return nil, me.Err(err, "failed to download url", &me.KV{"url", url})
	}
	sourceSize := len(bts)
	//resized, err := a.resize(bts, a.MaxHeight, a.MaxWidth)
	//if err != nil {
	//	return nil, me.Err(err, "failed to resize image", &me.KV{"url", url})
	//}

	bts = nil //release memory

	//	phash, w, h, err := a.phash(bytes.NewReader(resized))
	//	if err != nil {
	//		me.LogError(a.Log, "failed to generate phash on resized image", err)
	//	}

	resized, w, h, err := a.resize(bytes.NewReader(bts), a.MaxHeight, a.MaxWidth)

	uploaded, err := a.upload(resized, filename)
	if err != nil {
		return nil, me.Err(err, "failed up to upload resized image",
			&me.KV{"url", url},
			&me.KV{"filename", filename})
	}
	return &AcquiredImage{
		DestUrl:  uploaded.Url,
		DestSize: uploaded.FileSize,
		//	PHash:             phash,
		DestHeight:        h,
		DestWidth:         w,
		SourceUrl:         url,
		SourceContentType: ctype,
		SourceSize:        sourceSize,
	}, nil
}

func (a *AcquireImage) download(url string, timeout time.Duration) ([]byte, string, error) {
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, "", me.Err(err, "failed to get URL ")
	}
	defer resp.Body.Close()
	ct := strings.Join(resp.Header["Content-Type"], ",")
	bts, err := ioutil.ReadAll(resp.Body)
	return bts, ct, err
}

// phash - returns phash, w, h, err
//func (a *AcquireImage) phash(r io.Reader) (uint64, int, int error) {
//	img, err := magick.DecodeData(r)
//	if err != nil {
//		return 0, err
//	}
//	w := img.Width()
//	h := img.Height()
//	phash, err := img.Phash()
//	return phash, w, h, err
//}

func (a *AcquireImage) resizeVips(img []byte, maxWidth, maxHeight int) ([]byte, error) {
	options := vips.Options{
		Width:        maxWidth,
		Height:       maxHeight,
		Crop:         false,
		Enlarge:      false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}

	return vips.Resize(img, options)
}

func (a *AcquireImage) resize(imgR io.Reader, maxWidth, maxHeight int) (io.Reader, int, int, error) {

	img, err := imaging.Decode(imgR)
	if err != nil {
		return nil, 0, 0, err
	}

	thumb := imaging.Fit(img, maxWidth, maxHeight, imaging.CatmullRom)

	height := thumb.Rect.Dy()
	width := thumb.Rect.Dx()

	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, thumb, imaging.PNG)
	if err != nil {
		return nil, 0, 0, me.Err(err, "failed to encode")
	}
	return bytes.NewReader(buf.Bytes()), width, height, nil
}

func (a *AcquireImage) upload(source io.Reader, filename string) (*FileUploadOutput, error) {
	w, err := a.Uploader.Writer(filename)
	if err != nil {
		return nil, me.Err(err, "failed to get writer")
	}

	size, err := io.Copy(w, source)
	if err != nil {
		return nil, me.Err(err, "copy to upload fail")
	}

	if err = w.Close(); err != nil {
		return nil, me.Err(err, "close writer fail")
	}

	return &FileUploadOutput{
		Url:      a.Uploader.DestinationUrl(filename),
		FileSize: size,
	}, nil
}
