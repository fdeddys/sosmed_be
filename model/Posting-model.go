package model

import "time"

type Posting struct {
	Post       string    `json:"post"`
	Like       int32     `json:"like"`
	Dislike    int32     `json:"dislike"`
	Pict       string    `json:"pict"`
	InsertDate time.Time `json:"insertDate"`
	RestoId    int64     `json:"restoId"`
	Urut       int64     `json:"urut"`
}
