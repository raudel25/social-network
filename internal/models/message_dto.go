package models

type MessageRequest struct {
	RichText RichText `json:"richText" binding:"required"`
}
