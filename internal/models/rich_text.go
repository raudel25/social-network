package models

type RichText struct {
	Text string `json:"text" binding:"required"`
	HTML string `json:"html" binding:"required"`
}
