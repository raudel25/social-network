package models

type PostRequest struct {
	Title    string   `json:"title" binding:"required"`
	PhotoID  *uint    `json:"photo_id" binding:"required"`
	RePostID *uint    `json:"re_post_id" binding:"required"`
	RichText RichText `json:"rich_text" binding:"required"`
}
