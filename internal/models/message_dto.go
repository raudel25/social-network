package models

type MessageRequest struct {
	RichText RichText `json:"rich_text" binding:"required"`
}
