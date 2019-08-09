package model

import (
	"errors"
)

type Event struct {
	Bucket string `json:"bucket"`
	Folder string `json:"folder"`
	Width int     `json:"width"`
	Height int    `json:"height"`
}

func (event *Event) Validate() error {
	if event.Bucket == "" {
		return errors.New("'folder' has not been set in the event")
	}

	if event.Folder == "" {
		return errors.New("'folder' has not been set in the event")
	}

	if event.Width == 0 {
		return errors.New("'width' has not been set in the event")
	}

	if event.Height == 0 {
		return errors.New("'height' has not been set in the event")
	}

	return nil
}
