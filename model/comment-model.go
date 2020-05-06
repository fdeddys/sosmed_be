package model

import "time"

// Comment ...
type Comment struct {
	ID              string    `json:"id"`
	CommendParentID string    `json:"parentId"`
	AuthorName      string    `json:"authorName"`
	AuthorImageURL  string    `json:"authorImageURL"`
	Message         string    `json:"message"`
	InsertDate      time.Time `json:"timeComment"`
	PostID          string    `json:"postID"`
	Urut            int64     `json:"urut"`
}
