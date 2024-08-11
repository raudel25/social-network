package models

import "time"

type MessageRequest struct {
	RichText RichText `json:"richText" binding:"required"`
}

type MessageResponse struct {
	ID       uint            `json:"id"`
	Profile  ProfileResponse `json:"profile"`
	RichText RichText        `json:"richText"`
	Date     time.Time       `json:"date"`
}
