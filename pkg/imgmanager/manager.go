package imgmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

const MaxResizePower = 10

type ImgManager struct {
	path string
}

type Resizer interface {
	LoadAndResize(filename string, power uint8) (*ImageStat, error)
}

// New creates a new ImgManager
// path - is path to files
func New(path string) *ImgManager {
	return &ImgManager{
		path: path,
	}
}

// checkFolder checks the existence of the folder where files will be saved
// if there is no folder then the function will try a new folder otherwise an error will be returned
func (m *ImgManager) checkFolder() error {
	err := os.MkdirAll(m.path, 0755)
	if os.IsExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (m *ImgManager) ReadAndSaveNewImage(r io.Reader) (stat *ImageStat, err error) {
	img, ext, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	stat, err = m.save(img, generateFileName(ext))
	if err != nil {
		return nil, err
	}

	return stat, nil
}

// LoadAndResize loaded and resize image
// power means how much to divide the loaded image size
// for example if power = 2 then image width and height will be divided by power
// if power == 0 || power > MaxResizePower then power = MaxResizePower
func (m *ImgManager) LoadAndResize(filename string, power uint8) (*ImageStat, error) {
	if power == 0 || power > MaxResizePower {
		power = MaxResizePower
	}
	img, _, err := m.loadImage(filename)
	if err != nil {
		return nil, err
	}

	resizedImage := m.resize(img, power)
	newStat, err := m.save(resizedImage, filename)
	if err != nil {
		return nil, err
	}

	return newStat, nil
}

// resize is used to resize passed image.
func (m *ImgManager) resize(img image.Image, power uint8) image.Image {
	p := img.Bounds().Size()
	x := uint(p.X) / uint(power)
	y := uint(p.Y) / uint(power)

	return resize.Resize(x, y, img, resize.Lanczos3)
}

// save saves passed image in fs.
func (m *ImgManager) save(img image.Image, filename string) (stat *ImageStat, err error) {
	if err = m.checkFolder(); err != nil {
		return nil, err
	}
	f, err := os.Create(m.getFilePath(filename))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := encodeImage(f, img, filename); err != nil {
		return nil, err
	}

	return createImageStat(f, img, filename)
}

// getFilePath is used for getting file path based on file name.
func (m *ImgManager) getFilePath(filename string) string {
	return filepath.Join(m.path, filename)
}

// loadImage is used to load image from fs.
func (m *ImgManager) loadImage(filename string) (img image.Image, stat *ImageStat, err error) {
	f, err := os.Open(m.getFilePath(filename))
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	img, _, err = image.Decode(f)
	if err != nil {
		return nil, nil, err
	}

	stat, err = createImageStat(f, img, filename)
	if err != nil {
		return nil, nil, err
	}
	return img, stat, err
}

func encodeImage(w io.Writer, img image.Image, filename string) (err error) {
	switch filepath.Ext(filename) {
	case ".jpeg":
		return jpeg.Encode(w, img, nil)
	case ".png":
		return png.Encode(w, img)
	default:
		return image.ErrFormat
	}
}

// generateFileName generates filename based on uuid and file extension
func generateFileName(ext string) string {
	return uuid.NewString() + "." + ext
}

func createImageStat(f *os.File, img image.Image, filename string) (*ImageStat, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := img.Bounds().Size()
	return &ImageStat{
		ID:     filename,
		Width:  uint(size.X),
		Height: uint(size.Y),
		Size:   info.Size(),
	}, nil
}

type ImageStat struct {
	ID     string `json:"id"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Size   int64  `json:"size"`
}

func (m *ImageStat) MustMarshal() []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return b
}

func (m ImageStat) String() string {
	var kb float64
	if m.Size != 0 {
		kb = float64(m.Size) / 1024
	}

	return fmt.Sprintf("id=%s x=%d y=%d size=%.2fK", m.ID, m.Width, m.Height, kb)
}

func IsUnknowFormatErr(err error) bool {
	return errors.Is(err, image.ErrFormat)
}
