package database

import (
	"database/sql"
	"forum/internal/models"
)

type UserRepoInterface interface {
	CreateUserRepo(*models.User) (int64, error)
	GetUserByEmail(string) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	GetUserByUserID(int) (*models.User, error)
	CreateSession(*models.Session) error
	UpdateSession(*models.Session) error
	GetSessionByUserID(int) (*models.Session, error)
	GetSessionByToken(string) (*models.Session, error)
	DeleteSessionByToken(string) error
	DeleteSessionByUserID(int) error
	ChangeUserRole(string, int) error
	GetUserRole(int) (string, error)
	GetUserByRole(string) ([]*models.User, error)
}

type PostRepoInterface interface {
	CreatePostRepo(*models.Post) (int64, error)
	GetAllPosts() ([]*models.Post, error)
	GetCategoriesByPostID(int) ([]string, error)
	GetPostByID(int) (*models.Post, error)
	GetPostsByUserId(int) ([]*models.Post, error)
	GetPostsByLikes(int) ([]*models.Post, error)
	CreatePostCategory([]string, int) (int64, error)
	UpdateLikesCounter(int, int) error
	UpdateDislikesCounter(int, int) error
	GetReaction(int, int) (int, error)
	AddReactionToPostVotes(int, int, int) error
	DeleteFromPostVotes(int, int) error
	UpdateReactionInPostVotes(int, int, int) error
	GetPostsByCategory(string) ([]*models.Post, error)
}

type CommentRepoInterface interface {
	CreateCommentRepo(*models.Comment) (int64, error)
	GetAlCommentsForPost(int) ([]*models.Comment, error)
	UpdateLikesCounter(int, int) error
	UpdateDislikesCounter(int, int) error
	GetCommentReaction(int, int) int
	AddReactionToCommentVotes(int, int, int) error
	DeleteReactionFromCommentVotes(int, int) error
	UpdateReactionInCommentVotes(int, int, int) error
}

type Repository struct {
	UserRepoInterface
	PostRepoInterface
	CommentRepoInterface
}

func NewRepository(db *sql.DB) *Repository {
	repositoryObj := Repository{
		UserRepoInterface:    CreateNewUserDB(db),
		PostRepoInterface:    CreateNewPostDB(db),
		CommentRepoInterface: CreateNewCommentDB(db),
	}
	return &repositoryObj
}
