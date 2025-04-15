package exifdata

import (
	"fmt"
	"github.com/barasher/go-exiftool"
	"github.com/nobuenhombre/suikat/pkg/converter"
	"github.com/nobuenhombre/suikat/pkg/ge"
	"time"
)

type IExifData interface {
	GetDuration(fileName string) (time.Duration, error)
}

type conn struct{}

func New() IExifData {
	return &conn{}
}

func (c *conn) getFileExifData(fileName string) (exifData map[string]interface{}, err error) {
	exifData = make(map[string]interface{})

	et, err := exiftool.NewExiftool()
	if err != nil {
		return exifData, ge.Pin(err)
	}

	defer func() {
		etCloseErr := et.Close()
		if etCloseErr != nil {
			err = ge.Pin(etCloseErr, ge.Params{"err": err})
		}
	}()

	fileInfos := et.ExtractMetadata(fileName)
	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			return exifData, ge.Pin(fileInfo.Err)
		}

		exifData = fileInfo.Fields
	}

	return exifData, nil
}

func (c *conn) GetDuration(fileName string) (time.Duration, error) {
	exifData, err := c.getFileExifData(fileName)
	if err != nil {
		return time.Duration(0), ge.Pin(err)
	}

	for key, value := range exifData {
		if key == "Duration" {
			result, err := converter.StringToDuration(value.(string))
			if err != nil {
				return time.Duration(0), ge.Pin(err)
			}

			return result, nil
		}
	}

	return time.Duration(0), ge.Pin(fmt.Errorf("no Duration found"))
}
