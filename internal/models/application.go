package models

import (
	"database/sql"
	"time"
)

const (
	GoogleAuthURL      = "https://accounts.google.com/o/oauth2/auth"                                // const URL
	GoogleClientID     = "722031461724-dnvp1cl4hngcs1kgt0a2qi9j86a3dr1n.apps.googleusercontent.com" // my google account
	GoogleRedirectURL  = "https://localhost:8080/auth/google/callback"                              // callback endpoint
	GoogleClientSecret = "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"                                      // my google account
)

const (
	GitHubAuthURL      = "https://github.com/login/oauth/authorize"
	GitHubClientID     = "7204d1f96b4db7e5d453"
	GitHubRedirectURL  = "https://localhost:8080/auth/github/callback"
	GitHubClientSecret = "2a4621170475143853a9752bf405fb1d2f781051"
)

type GitHubResponseToken struct {
	AccessToken string `json:"access_token"`
	TokenID     string `json:"id_token"`
	Scope       string `json:"scope"`
}

type GoogleResponseToken struct {
	AccessToken string `json:"access_token"`
	TokenID     string `json:"id_token"`
}

type GoogleUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
	Password       string
}

type GitHubUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
	Password       string
}

type GoogleLoginUserData struct {
	ID         int
	Name       string
	Email      string
	Password   string
	FirstName  string
	SecondName string
	Provider   string
}

type GitHubLoginUserData struct {
	ID         int
	Name       string
	Email      string
	Password   string
	FirstName  string
	SecondName string
	Login      string
	Provider   string
}

type User struct {
	UserID     int
	FirstName  string
	SecondName string
	Username   string
	Email      string
	Password   string
	Role       string
}

type Session struct {
	UserID  int
	Token   string
	ExpTime time.Time
}

type Post struct {
	PostID            int
	UserID            int
	Username          string
	Title             string
	Content           string
	CreatedTime       time.Time
	CreatedTimeString string
	LikesCounter      int
	DislikeCounter    int
	Categories        []string
	Comments          []*Comment
	ImagePath         string
	UserRole          string
}

type Comment struct {
	CommentID         int
	PostID            int
	UserID            int
	Username          string
	Content           string
	CreatedTime       time.Time
	CreatedTimeString string
	LikesCounter      int
	DislikeCounter    int
	UserRole          string
}

type Database struct {
	DB *sql.DB
}
