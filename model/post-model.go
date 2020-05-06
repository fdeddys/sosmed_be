package model

import "time"

// Post ...
type Post struct {
	ID           string    `json:"id"`
	Post         string    `json:"posting"`
	Like         int32     `json:"like"`
	Dislike      int32     `json:"dislike"`
	Pict         string    `json:"imgUrl"`
	InsertDate   time.Time `json:"timePost"`
	RestoID      int64     `json:"restoId"`
	RestoName    string    `json:"restoName"`
	RestoImgURL  string    `json:"restoImageUrl"`
	Urut         int64     `json:"urut"`
	TotalComment int64     `json:"totalComment"`
}
