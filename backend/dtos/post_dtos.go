package dtos

import "time"

type AuthorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PostResponse struct {
	ID         uint               `json:"id"`
	Title      string             `json:"title"`
	Content    string             `json:"content"`
	Image      string             `json:"image"`
	CreatedAt  time.Time          `json:"created_at"`
	Author     AuthorResponse     `json:"author"`
	Categories []CategoryResponse `json:"categories"`
	Likes      int                `json:"likes"`
	Dislikes   int                `json:"dislikes"`
}

type CreatePostRequest struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Image       string `json:"image"`
	CategoryIDs []uint `json:"category_ids"`
}
