package models

import "time"

type PostRequest struct {
	Title    string   `json:"title" binding:"required"`
	PhotoID  *uint    `json:"photo_id"`
	RePostID *uint    `json:"re_post_id"`
	RichText RichText `json:"rich_text" binding:"required"`
}

type PostResponse struct {
	Title         string          `json:"title"`
	Profile       ProfileResponse `json:"profile"`
	Photo         *Photo          `json:"photo"`
	RichText      RichText        `json:"rich_text"`
	RePost        *PostResponse   `json:"re_post"`
	Reaction      bool            `json:"reaction"`
	CantReactions int             `json:"cant_reactions"`
	CantMessages  int             `json:"cant_messages"`
	Date          time.Time       `json:"date"`
}
