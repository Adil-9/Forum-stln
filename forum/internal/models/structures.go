package models

import "time"

type Comment struct {
	ID           int
	UserID       int
	PostID       int
	LikeCount    int
	DislikeCount int
	Vote         int
	Content      string
	Author       string
}

type Post struct {
	ID           int
	AuthorID     int
	LikeCount    int
	DislikeCount int
	CommentCount int
	Vote         int
	Author       string
	Title        string
	Content      string
	Categories   []string
}

type Reaction struct {
	ID        string
	UserID    int
	PostID    int
	CommentID int
	Vote      int
}

type Session struct {
	ID             int
	UserID         int
	Token          string
	ExpirationDate time.Time
}

type TemplateData struct {
	Template string
	User     User
	Post     Post
	Posts    []Post
	Comments []Comment
	Error    ErrorMsg
}

type ErrorMsg struct {
	Status int
	Msg    string
}

type User struct {
	ID              int
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}
