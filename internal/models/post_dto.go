package models

import "time"

type PostRequest struct {
	PhotoID  *uint     `json:"photoId"`
	RePostID *uint     `json:"rePostId"`
	RichText *RichText `json:"richText" binding:"required"`
}

type PostResponse struct {
	ID            uint              `json:"id"`
	Profile       ProfileResponse   `json:"profile"`
	PhotoID       *uint             `json:"photoId"`
	RichText      *RichText         `json:"richText"`
	RePost        *PostResponse     `json:"rePost"`
	Reaction      bool              `json:"reaction"`
	CantReactions int               `json:"cantReactions"`
	CantMessages  int               `json:"cantMessages"`
	CantRePosts   int               `json:"cantRePosts"`
	Messages      []MessageResponse `json:"messages"`
	Date          time.Time         `json:"date"`
}
