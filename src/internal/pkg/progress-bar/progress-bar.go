package progressbar

import (
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type IProgress interface {
	Set(value int) error
}

func New() IProgress {
	return progressbar.NewOptions(100,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("[cyan]FFMPEG... [reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]━[reset]",
			SaucerHead:    "[red]━[reset]",
			SaucerPadding: "[light_gray]━[reset]",
			BarStart:      "",
			BarEnd:        "",
		}))
}
