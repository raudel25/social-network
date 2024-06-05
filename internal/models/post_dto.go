package models

import "time"

type PostRequest struct {
	PhotoID  *uint    `json:"photoId"`
	RePostID *uint    `json:"rePostId"`
	RichText *RichText `json:"richText" binding:"required"`
}

type PostResponse struct {
	Profile       ProfileResponse `json:"profile"`
	PhotoID       uint          `json:"photoId"`
	RichText      *RichText        `json:"richText"`
	RePost        *PostResponse   `json:"rePost"`
	Reaction      bool            `json:"reaction"`
	CantReactions int             `json:"cantReactions"`
	CantMessages  int             `json:"cantMessages"`
	Date          time.Time       `json:"date"`
}
