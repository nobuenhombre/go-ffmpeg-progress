package progresspipe

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nobuenhombre/suikat/pkg/converter"
	"github.com/nobuenhombre/suikat/pkg/fifo"
	"github.com/nobuenhombre/suikat/pkg/ge"
	progressbar "go-ffmpeg-progress/src/internal/pkg/progress-bar"
	"math"
	"strings"
	"time"
)

type IProgressPipe interface {
	GetConverterArgs() []string
	ReadProgress(duration time.Duration) error
}

type conn struct {
	Name        string
	Pipe        fifo.Service
	ProgressBar progressbar.IProgress
}

func New(progress progressbar.IProgress) (IProgressPipe, error) {
	pipeUUID := uuid.New()
	name := fmt.Sprintf("/tmp/ffmpeg-progress-pipe-%v", pipeUUID.String())

	pipe := fifo.New(name)

	err := pipe.Create()
	if err != nil {
		return nil, ge.Pin(err)
	}

	return &conn{
		Name:        name,
		Pipe:        pipe,
		ProgressBar: progress,
	}, nil
}

func (c *conn) GetConverterArgs() []string {
	return []string{
		"-progress",
		c.Name,
	}
}

func (c *conn) getReceiver(duration time.Duration) fifo.MessageReceiver {
	return func(msg string) error {
		if strings.Contains(msg, "out_time_us=N/A") {
			return nil
		}

		if strings.Contains(msg, "out_time_us=") {
			currentDurationMs, err := converter.StringToInt64(
				strings.TrimSpace(strings.ReplaceAll(msg, "out_time_us=", "")),
			)
			if err != nil {
				return ge.Pin(err)
			}

			progress := int(math.Round(float64(time.Duration(currentDurationMs)*time.Microsecond*100) / float64(duration)))

			err = c.ProgressBar.Set(progress)
			if err != nil {
				return ge.Pin(err)
			}
		}

		return nil
	}
}

func (c *conn) ReadProgress(duration time.Duration) error {
	rcv := c.getReceiver(duration)

	err := c.Pipe.OpenToRead()
	if err != nil {
		return ge.Pin(err)
	}

	defer func() {
		fifoCloseErr := c.Pipe.Close()
		if fifoCloseErr != nil {
			err = ge.Pin(fifoCloseErr, ge.Params{"err": err})
		}

		fifoDeleteErr := c.Pipe.Delete()
		if fifoDeleteErr != nil {
			err = ge.Pin(fifoDeleteErr, ge.Params{"err": err})
		}
	}()

	err = c.Pipe.Read(rcv)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}
