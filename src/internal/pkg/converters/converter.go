package converters

import (
	"github.com/nobuenhombre/suikat/pkg/ge"
	"github.com/nobuenhombre/suikat/pkg/osexec"
)

type IConverter interface {
	Convert(args []string) error
}

type Conn struct {
	Cmd string
}

func New(cmd string) IConverter {
	return &Conn{
		Cmd: cmd,
	}
}

func (c *Conn) Convert(args []string) error {
	_, err := osexec.OSRun(c.Cmd, args)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}
