package imgmanager

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

var ErrFormat = image.ErrFormat

type ImgManager struct {
	path string
}

// New creates a new ImgManager
// path - is path to files
func New(path string) *ImgManager {
	return &ImgManager{
		path: path,
	}
}

func (m *ImgManager) SaveImage(r io.Reader) (filename string, err error) {
	img, ext, err := decodeImage(r)
	if err != nil {
		return "", err
	}
	filename = generateFileName(ext)
	filepath := m.getFilePath(filename)
	f, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = encodeImage(f, img, ext)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (m *ImgManager) getFilePath(filename string) string {
	return filepath.Join(m.path, filename)
}

func (m *ImgManager) DoResize(filename string) error {
	filepath := m.getFilePath(filename)
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	img, ext, err := decodeImage(f)
	f.Close()
	if err != nil {
		return err
	}
	f, err = os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	img = resize.Resize(1000, 0, img, resize.Lanczos3)
	return encodeImage(f, img, ext)
}

func encodeImage(w io.Writer, img image.Image, ext string) (err error) {
	switch ext {
	case "jpeg":
		err = jpeg.Encode(w, img, nil)
	case "png":
		err = png.Encode(w, img)
	default:
		err = ErrFormat
	}
	return err
}

func decodeImage(r io.Reader) (image.Image, string, error) {
	img, ext, err := image.Decode(r)
	if err == image.ErrFormat {
		return nil, "", ErrFormat
	}
	if err != nil {
		return nil, "", err
	}
	return img, ext, nil
}

func generateFileName(ext string) string {
	return uuid.NewString() + "." + ext
}
