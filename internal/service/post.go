package service

import (
	"errors"
	"fmt"
	"forum/internal/database"
	"forum/internal/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PostServiceImpl struct {
	repo database.PostRepoInterface
}

func CreateNewPostService(repo database.PostRepoInterface) *PostServiceImpl {
	postService := PostServiceImpl{repo: repo}
	return &postService
}

func (postObj *PostServiceImpl) CreatePost(post *models.Post, userRole string) (int, int, error) {
	if err := postObj.isPostParamsValid(post); err != nil {
		return http.StatusBadRequest, -1, err
	}

	post.CreatedTime = time.Now()
	post.LikesCounter = 0
	post.DislikeCounter = 0

	if userRole == "admin" || userRole == "moderator" {
		post.IsApproved = 1
	} else {
		post.IsApproved = 0
	}
	id, err := postObj.repo.CreatePostRepo(post)
	if err != nil {
		return http.StatusInternalServerError, -1, err
	}

	post.PostID = int(id)

	_, err = postObj.createPostCategory(post.Categories, post.PostID)
	if err != nil {
		return http.StatusInternalServerError, -1, err
	}
	return http.StatusOK, int(id), nil
}

func (postObj *PostServiceImpl) GetAllPosts() ([]*models.Post, error) {
	posts, err := postObj.repo.GetAllPosts()
	// fmt.Println("fetched all postss")
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (postObg *PostServiceImpl) isPostParamsValid(post *models.Post) error {
	if len(post.Title) < 2 {
		return errors.New("The title must be at least 2 characters")
	}
	if len(post.Content) < 2 {
		return errors.New("The content must be at least 2 characters")
	}
	if len(post.Categories) == 0 {
		return errors.New("Didn't select the categories you want")
	}
	return nil
}

func (postObj *PostServiceImpl) GetCategories(postID int) ([]string, error) {
	categories, err := postObj.repo.GetCategoriesByPostID(postID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (postObj *PostServiceImpl) createPostCategory(categories []string, postID int) (int, error) {
	if err := postObj.isCategoryValid(categories); err != nil {
		return -1, err
	}

	id, err := postObj.repo.CreatePostCategory(categories, postID)
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func (postObj *PostServiceImpl) isCategoryValid(categories []string) error {
	if len(categories) == 0 {
		return errors.New("YOUR CATEGORY IS NULL")
	}
	return nil
}

func (postObj *PostServiceImpl) UpdateReaction(currReaction int, postID int, userID int) error {
	var err error
	prevReaction, _ := postObj.repo.GetReaction(postID, userID)

	if prevReaction == 0 {
		postObj.repo.AddReactionToPostVotes(postID, userID, currReaction)
		if currReaction == 1 {
			err = postObj.repo.UpdateLikesCounter(postID, 1)
		} else {
			err = postObj.repo.UpdateDislikesCounter(postID, 1)
		}
	} else if prevReaction == currReaction {
		postObj.repo.DeleteFromPostVotes(postID, userID)
		if currReaction == 1 {
			postObj.repo.UpdateLikesCounter(postID, -1)
		} else {
			postObj.repo.UpdateDislikesCounter(postID, -1)
		}
	} else if prevReaction != currReaction {
		postObj.repo.UpdateReactionInPostVotes(postID, userID, currReaction)
		if currReaction == 1 {
			postObj.repo.UpdateLikesCounter(postID, 1)
			postObj.repo.UpdateDislikesCounter(postID, -1)
		} else {
			postObj.repo.UpdateLikesCounter(postID, -1)
			postObj.repo.UpdateDislikesCounter(postID, 1)
		}
	}
	return err
}

func (postObj *PostServiceImpl) GetPostByID(postID int) (*models.Post, error) {
	post, err := postObj.repo.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (postObj *PostServiceImpl) Filter(field string, userID int) ([]*models.Post, error) {
	posts := []*models.Post{}
	var err error

	if field == "CreatedPosts" {
		posts, err = postObj.repo.GetPostsByUserId(userID)
		if err != nil {
			return nil, err
		}
	} else if field == "LikedPosts" {
		posts, err = postObj.repo.GetPostsByLikes(userID)
		if err != nil {
			return nil, err
		}
	} else {
		field = strings.ToLower(field)
		posts, err = postObj.repo.GetPostsByCategory(field)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (postObj *PostServiceImpl) AddImagesToPost(files *multipart.FileHeader) (string, error) {
	openedfile, err := files.Open()
	if err != nil {
		return "Failed: Couldn't add image", err
	}
	defer openedfile.Close()

	buff := make([]byte, 512)
	_, err = openedfile.Read(buff)
	if err != nil {
		return "Failed in reading the file", err
	}

	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/jpg" && fileType != "image/gif" {
		return "Failed: Image has inappropriate format", errors.New("Invalid File Type")
	}

	_, err = openedfile.Seek(0, io.SeekStart)
	if err != nil {
		return "Image cannot be shown", errors.New("Invalid File Type")
	}

	// to keep the images visible we need this or to set the frame?
	err = os.MkdirAll("./data/assets/images", os.ModePerm)
	if err != nil {
		return "Failed: coudn't create the directory", err
	}

	imageName := uuid.New().String()
	imageDest := fmt.Sprintf("/images/%d%s", imageName, filepath.Ext(files.Filename))

	// creates a new file at the destination path. If an error occurs, it returns the error.
	dst, err := os.Create("./data/assets" + imageDest)
	if err != nil {
		return "Failed: Could not create image in the directory", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, openedfile)
	if err != nil {
		return "Failed: couldn't copy image to the directory", err
	}
	return imageDest, nil
}

func (postObj *PostServiceImpl) DeletePost(postID int) error {
	err := postObj.repo.DeletePostByID(postID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostServiceImpl) DeletePostCategoryByPostID(postID int) error {
	err := postObj.repo.DeletePostCategoryByPostID(postID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostServiceImpl) DeleteAllPostVotesByPostID(postID int) error {
	err := postObj.repo.DeleteAllPostVotesByPostID(postID)
	if err != nil {
		return err
	}
	return nil
}

func (postObj *PostServiceImpl) ApprovePost(postID int) error {
	err := postObj.repo.UpdateIsApprovePostStatus(postID)
	if err != nil {
		return err
	}
	return nil
}
