package common

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateChecksum(t *testing.T) {
	tests := []struct {
		name             string
		filePath         string
		expected         string
		isError          bool
		content          []byte
		isMultimediaFile bool
	}{
		{
			name:     "empty file",
			filePath: "empty.txt",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			content:  []byte(""),
			isError:  false,
		},
		{
			name:     "small file",
			filePath: "small.txt",
			expected: "4b4f9dde5da9baa5184b3872440ac93dfa77ab365476e5317257eee49f89d09c",
			content:  []byte("This is a small test file"),
			isError:  false,
		},
		{
			name:     "large file",
			filePath: "large.txt",
			expected: "e5b844cc57f57094ea4585e235f36c78c1cd222262bb89d53c94dcb4d6b3e55d",
			content:  make([]byte, (1<<20)*10),
			isError:  false,
		},
		{
			name:             "multimedia file",
			filePath:         "rgb.png",
			expected:         "e5b844cc57f57094ea4585e235f36c78c1cd222262bb89d53c94dcb4d6b3e55d",
			isError:          false,
			isMultimediaFile: true,
		},
		{
			name:     "non-existent file",
			filePath: "non-existent.txt",
			expected: "",
			isError:  true,
		},
		{
			name:     "directory",
			filePath: "test-data",
			expected: "",
			isError:  true,
		},
	}

	for _, tt := range tests {
		var f *os.File
		var err error
		if !tt.isError {
			testFileDir, err := os.MkdirTemp("", "testfile")
			if err != nil {
				t.Fatalf("Create temporary test dir error: %v", err)
			}

			filePath := filepath.Join(testFileDir, tt.filePath)
			f, err := os.Create(filePath)
			if err != nil {
				t.Fatalf("Create temporary test file error: %v", err)
			}
			defer os.Remove(f.Name())
			defer f.Close()

			if tt.isMultimediaFile {
				m := generateImage()
				png.Encode(f, m)
			} else {
				if _, err = f.Write(tt.content); err != nil {
					t.Fatalf("Write content to temporary test file error: %v", err)
				}
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if tt.isError {
				_, err = GenerateChecksum(tt.filePath)
				if err == nil {
					t.Errorf("Expected error, but got none")
				}
			} else {
				actual, err := GenerateChecksum(f.Name())
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if actual != tt.expected {
					t.Errorf("Expected: %v, but got: %v", tt.expected, actual)
				}
			}
		})
	}
}

type Circle struct {
	X, Y, R float64
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	} else {
		return 255
	}
}

func generateImage() *image.RGBA {
	var w, h int = 280, 240
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := 40.0
	θ := 2 * math.Pi / 3
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 60}
	cg := &Circle{hw - r*math.Sin(θ), hh - r*math.Cos(θ), 60}
	cb := &Circle{hw - r*math.Sin(-θ), hh - r*math.Cos(-θ), 60}

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.RGBA{
				cr.Brightness(float64(x), float64(y)),
				cg.Brightness(float64(x), float64(y)),
				cb.Brightness(float64(x), float64(y)),
				255,
			}
			m.Set(x, y, c)
		}
	}
	return m
}
