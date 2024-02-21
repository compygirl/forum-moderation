package handlers

import (
	"errors"
	"forum/internal/models"
	helpers "forum/internal/web/handlers/helpers"
	"net/http"
)

func (h *Handler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		helpers.ErrorHandler(w, http.StatusNotFound, errors.New(" "))
		return
	}

	type templateData struct {
		LoggedIn      bool
		AllPosts      []*models.Post
		AllCategories []string
		// Role          string
	}

	posts, err := h.service.PostServiceInterface.GetAllPosts()
	if err != nil {
		helpers.ErrorHandler(w, http.StatusInternalServerError, err)
		return
	}
	for _, post := range posts {
		// getting all categories
		temp_categories, err := h.service.PostServiceInterface.GetCategories(post.PostID)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}
		post.Categories = append(post.Categories, temp_categories...)

		// getting the username for posts
		user, err := h.service.UserServiceInterface.GetUserByUserID(post.UserID)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}
		post.Username = user.Username

		// post.ImagePath = path

		// changing the format of the time
		post.CreatedTimeString = post.CreatedTime.Format("Jan 2, 2006 at 15:04")

	}

	indexPath := "internal/web/templates/index.html"
	data := templateData{
		LoggedIn:      h.service.IsUserLoggedIn(r),
		AllPosts:      posts,
		AllCategories: []string{"Movie", "Game", "Book", "Others"}, // Initialize AllCategories with values
		// Role:          user.Role,
	}
	helpers.RenderTemplate(w, indexPath, data)
}
