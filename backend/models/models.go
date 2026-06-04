package models

import (
	"gorm.io/gorm"
	"time"
)

type ReactionType string

const (
	Like    ReactionType = "like"
	Dislike ReactionType = "dislike"
)

type User struct {
	gorm.Model
	UserName      string         `gorm:"column:username;not null;unique" json:"username"`
	Password      string         `gorm:"column:password;not null"        json:"password"`
	Email         string         `gorm:"column:email;not null;unique"    json:"email"`
	Avatar        string         `gorm:"column:avatar"                   json:"avatar"`
	Posts         []Post         `gorm:"foreignKey:UserID"`
	Comments      []Comment      `gorm:"foreignKey:UserID"`
	Reactions     []Reaction     `gorm:"foreignKey:UserID"`
	RefreshTokens []RefreshToken `gorm:"foreignKey:UserID"`
}

type PasswordReset struct {
	gorm.Model
	UserID    uint      `gorm:"column:user_id;not null" json:"userId"`
	Token     string    `gorm:"column:token;not null"   json:"token"`
	ExpiresAt time.Time `gorm:"column:expires_at"       json:"expiresAt"`
	Used      bool      `gorm:"column:used;default:false" json:"used"`
	User      User      `gorm:"foreignKey:UserID"`
}

type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"column:user_id;not null" json:"userId"`
	Token     string    `gorm:"column:token;not null"   json:"token"`
	ExpiresAt time.Time `gorm:"column:expires_at"       json:"expiresAt"`
	Revoked   bool      `gorm:"column:revoked;default:false" json:"revoked"`
	User      User      `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	UserID     uint       `gorm:"column:user_id;not null" json:"userId"`
	Title      string     `gorm:"column:title;not null"   json:"title"`
	Content    string     `gorm:"column:content;not null" json:"content"`
	Image      string     `gorm:"column:image"            json:"image"`
	User       User       `gorm:"foreignKey:UserID"`
	Comments   []Comment  `gorm:"foreignKey:PostID"`
	Reactions  []Reaction `gorm:"foreignKey:PostID"`
	Categories []Category `gorm:"many2many:post_categories"`
}

type Comment struct {
	gorm.Model
	UserID  uint   `gorm:"column:user_id;not null" json:"userId"`
	PostID  uint   `gorm:"column:post_id;not null" json:"postId"`
	Content string `gorm:"column:content;not null" json:"content"`
	User    User   `gorm:"foreignKey:UserID"`
	Post    Post   `gorm:"foreignKey:PostID"`
}

type Reaction struct {
	gorm.Model
	UserID uint         `gorm:"column:user_id;not null;uniqueIndex:idx_user_post_reaction" json:"userId"`
	PostID uint         `gorm:"column:post_id;not null;uniqueIndex:idx_user_post_reaction" json:"postId"`
	Type   ReactionType `gorm:"column:type;not null"                                       json:"type"`
	User   User         `gorm:"foreignKey:UserID"`
	Post   Post         `gorm:"foreignKey:PostID"`
}

type Category struct {
	gorm.Model
	Name  string `gorm:"column:name;not null;unique" json:"name"`
	Posts []Post `gorm:"many2many:post_categories"`
}