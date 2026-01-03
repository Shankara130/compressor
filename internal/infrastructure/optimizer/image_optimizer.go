package optimizer

import (
	"image/jpeg"
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

	img, err := jpeg.Decode(in)
	if err != nil {
		return err
	}

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: 75})
}
