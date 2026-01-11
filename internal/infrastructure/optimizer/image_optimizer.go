package optimizer

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

type ImageOptimizer struct{}

func NewImageOptimizer() *ImageOptimizer {
	return &ImageOptimizer{}
}

func (o *ImageOptimizer) Optimize(input string, output string) error {
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()

	img, format, err := image.Decode(in)
	if err != nil {
		return err
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 75})
	case "png":
		return png.Encode(out, img)
	default:
		return nil
	}

}
