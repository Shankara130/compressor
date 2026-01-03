package optimizer

import "os/exec"

type PDFOptimizer struct{}

func NewPDFOptimizer() *PDFOptimizer {
	return &PDFOptimizer{}
}

func (o *PDFOptimizer) Optimize(input, output string) error {
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
