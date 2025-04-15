package domainapp

import (
	"github.com/nobuenhombre/suikat/pkg/ge"
	"go-ffmpeg-progress/src/internal/pkg/converters"
	exifdata "go-ffmpeg-progress/src/internal/pkg/exif-data"
	progressbar "go-ffmpeg-progress/src/internal/pkg/progress-bar"
	progresspipe "go-ffmpeg-progress/src/internal/pkg/progress-pipe"
	"log"
	"os"
)

type AppDomain struct {
	ExifData        exifdata.IExifData
	ProgressBar     progressbar.IProgress
	ProgressPipe    progresspipe.IProgressPipe
	FFMPEGConverter converters.IConverter
}

func New() (IDomainApp, error) {
	exifData := exifdata.New()

	progressBar := progressbar.New()

	progressPipe, err := progresspipe.New(progressBar)
	if err != nil {
		return nil, ge.Pin(err)
	}

	ffmpegConverter := converters.New("/usr/bin/ffmpeg")

	return &AppDomain{
		ExifData:        exifData,
		ProgressBar:     progressBar,
		ProgressPipe:    progressPipe,
		FFMPEGConverter: ffmpegConverter,
	}, err
}

func (d *AppDomain) GetArgs() []string {
	return os.Args[1:]
}

func (d *AppDomain) GetInputFileName() (string, error) {
	args := d.GetArgs()
	for i, arg := range args {
		if arg == "-i" && len(args) > i+1 {
			return args[i+1], nil
		}
	}

	return "", ge.Pin(&ge.NotFoundError{
		Key: "key input file",
	})
}

func (d *AppDomain) Run() (err error) {
	inputFileName, err := d.GetInputFileName()
	if err != nil {
		return ge.Pin(err)
	}

	duration, err := d.ExifData.GetDuration(inputFileName)
	if err != nil {
		return ge.Pin(err)
	}

	go func() {
		err = d.ProgressPipe.ReadProgress(duration)
		if err != nil {
			log.Fatal(err)
		}
	}()

	errRun := d.FFMPEGConverter.Convert(
		append(
			d.GetArgs(),
			d.ProgressPipe.GetConverterArgs()...,
		),
	)
	if errRun != nil {
		return ge.Pin(err)
	}

	return nil
}
