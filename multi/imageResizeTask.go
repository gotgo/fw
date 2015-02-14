package multi

import (
	"os"
	"path/filepath"

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
