package multi

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/amattn/deeperror"
	"github.com/disintegration/imaging"
	"github.com/gotgo/fw/me"
)

type ImageResizeOutput struct {
	FileSize    int64
	Height      int
	Width       int
	ContentType string
	FilePath    string
}

type ResizeImageTask struct {
	MaxWidth  int
	MaxHeight int
}

func (r *ResizeImageTask) Run(input interface{}) (interface{}, error) {
	filepath := input.(string)
	return resize(filepath, r.MaxWidth, r.MaxHeight)
}

func (r *ResizeImageTask) Name() string {
	return "resizeImage"
}

func resize(filePath string, maxWidth, maxHeight int) (*ImageResizeOutput, error) {
	fmt.Printf("opening file to resize " + filePath)
	img, err := imaging.Open(filePath)
	if err != nil {
		return nil, me.Err(err, "failed to open "+filePath)
	}

	thumb := imaging.Fit(img, maxWidth, maxHeight, imaging.CatmullRom)

	height := thumb.Rect.Dy()
	width := thumb.Rect.Dx()

	dir, file := filepath.Split(filePath)
	ext := filepath.Ext(file)
	thumbnailPath := filepath.Join(dir, file[0:len(file)-len(ext)])
	f, err := os.Create(thumbnailPath)
	if err != nil {
		return nil, me.Err(err, "failed to create new emtpy file")
	}
	err = imaging.Encode(f, thumb, imaging.PNG)
	if err != nil {
		return nil, me.Err(err, "failed to encode")
	}

	err = f.Sync()
	if err != nil {
		return nil, me.Err(err, "failed to flush to disk")
	}

	var size int64
	if fi, err := f.Stat(); err == nil {
		size = fi.Size()
	}

	err = f.Close()
	if err != nil {
		return nil, me.Err(err, "failed to close")
	}

	// save the combined image to file
	//	err = imaging.Save(thumb, thumbnailPath)
	//if err != nil {
	//	return "", me.Err(err, "failed to save")
	//}
	return &ImageResizeOutput{
		FilePath:    thumbnailPath,
		Height:      height,
		Width:       width,
		ContentType: "image/png",
		FileSize:    size,
	}, nil
}

func getContentType(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", deeperror.New(rand.Int63(), "failed to open image", err)
	}
	defer file.Close()
	//if fi, err := file.Stat(); err == nil {
	//	i.Term.Size = fi.Size()
	//}
	rdr := bufio.NewReader(file)
	bts, _ := rdr.Peek(512)
	contentType := http.DetectContentType(bts)

	//m, _, err := image.Decode(rdr)
	_, _, err = image.Decode(rdr)
	if err != nil {
		return contentType, me.Err(err, "Failed to decode image")
	}

	return contentType, nil

	//	if m != nil {
	//		size := m.Bounds().Size()
	//		i.Term.Width = size.X
	//		i.Term.Height = size.Y
	//		i.Log.Debug(fmt.Sprintf("imaged type %s width:%d height:%d", contentType, size.X, size.Y))
	//	}
	//	return nil
}
