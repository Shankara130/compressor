package optimizer

import (
	"os/exec"

	"github.com/Shankara130/compressor/internal/domain/service"
)

type PDFOptimizer struct{}

func NewPDFOptimizer() *PDFOptimizer {
	return &PDFOptimizer{}
}

func (o *PDFOptimizer) Optimize(input, output string, _ service.ProgressFunc) error {
	cmd := exec.Command(
		"gs",
		"-sDEVICE=pdfwrite",
		"-dPDFSETTINGS=/ebook",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile="+output,
		input,
	)
	return cmd.Run()
}
