package imgmanager

import (
	"errors"
	"image"
	"testing"
)

func TestImgManager(t *testing.T) {
	m := New("testsdata/images")

	t.Run("testGetFilePath", testGetFilePath(m))
}

func testGetFilePath(m *ImgManager) func(t *testing.T) {
	return func(t *testing.T) {
		tests := []struct {
			name string
			arg  string
			want string
		}{
			{
				"ok",
				"foo.png",
				m.path + "/foo.png",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := m.getFilePath(tt.arg); got != tt.want {
					t.Errorf("getFilePath(%s) = %v, want %v", tt.arg, got, tt.want)
				}
			})
		}
	}
}

func TestIsUnknowFormatErr(t *testing.T) {
	tests := []struct {
		name string
		arg  error
		want bool
	}{
		{
			"ok",
			image.ErrFormat,
			true,
		},
		{
			"foo error",
			errors.New("foo error"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnknowFormatErr(tt.arg); got != tt.want {
				t.Errorf("IsUnknowFormatErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateFileName(t *testing.T) {
	var ext = ".jpeg"
	filename := generateFileName(ext)
	filename2 := generateFileName(ext)
	if filename == filename2 {
		t.Errorf(`Want unique filename %s, got %s`, filename, filename2)
	}
}
