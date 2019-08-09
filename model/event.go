package model

type Event struct {
	Bucket string `json:"bucket"`
	Folder string `json:"folder"`
	Width int     `json:"width"`
	Height int    `json:"height"`
}