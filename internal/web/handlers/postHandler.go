package handlers

import (
	"errors"
	"forum/internal/models"
	helpers "forum/internal/web/handlers/helpers"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		cookie := helpers.SessionCookieGet(r)
		const MaxImageSize = 20 * 1024 * 1024
		if cookie == nil {
			helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Conversion of postID failed"))
			return
		}
		session, err := h.service.UserServiceInterface.GetSession(cookie.Value)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Conversion of postID failed"))
			return
		}

		post := &models.Post{
			UserID:     session.UserID,
			Title:      r.FormValue("posttitle"),
			Content:    r.FormValue("postcontent"),
			Categories: r.Form["preference"],
		}

		//=============================================================
		//block of code responsible for the image upload
		r.Body = http.MaxBytesReader(w, r.Body, MaxImageSize)

		err = r.ParseMultipartForm(MaxImageSize)

		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Image size over 20 Mb"))
			return
		}

		if len(r.MultipartForm.File["files"]) != 0 {
			r.ParseForm()
			if r.MultipartForm.File["files"][0].Size > int64(MaxImageSize) {
				helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Image size over 20 Mb"))
			}
			file := r.MultipartForm.File["files"][0] // since only one image at a time
			path, err := h.service.AddImagesToPost(file)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, err)
				return
			}
			post.ImagePath = path
		}
		//=============================================================
		statusCode, postId, err := h.service.PostServiceInterface.CreatePost(post)
		post.PostID = postId
		if err != nil {
			helpers.ErrorHandler(w, statusCode, err)
			return
		}

		expTime, err := h.service.UserServiceInterface.ExtendSessionTimeout(cookie.Value)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("The Time cannot be extended"))
			return
		}
		err = helpers.SessionCookieExtend(r, w, expTime)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("The Time cannot be extended"))
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		helpers.ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("Error in Post Handler"))
		return

	}
}

func (h *Handler) ReactOnPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			helpers.ErrorHandler(w, http.StatusBadRequest, errors.New("Conversion of postID failed"))
			return
		}

		currReaction, err := strconv.Atoi(r.FormValue("type"))
		if err != nil {
			helpers.ErrorHandler(w, http.StatusBadRequest, errors.New("Conversion of reaction type failed"))
			return
		}

		cookie := helpers.SessionCookieGet(r)
		if cookie == nil {
			helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Cookie cannot be reseived in Post Reaction Handler"))
			return
		}

		session, err := h.service.UserServiceInterface.GetSession(cookie.Value)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}
		expTime, err := h.service.UserServiceInterface.ExtendSessionTimeout(cookie.Value)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Cookie cannot be extended"))
			return
		}
		err = helpers.SessionCookieExtend(r, w, expTime)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}

		if err := h.service.PostServiceInterface.UpdateReaction(currReaction, postID, session.UserID); err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	default:
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Error in Post Reaction Handler"))
		return
	}
}

func (h *Handler) FilterHandler(w http.ResponseWriter, r *http.Request) {
	type templateData struct {
		LoggedIn      bool
		AllPosts      []*models.Post
		AllCategories []string
	}

	switch r.Method {
	case "GET":
		var userID int
		cookie := helpers.SessionCookieGet(r)
		if cookie == nil {
			userID = 0
		} else {
			session, err := h.service.UserServiceInterface.GetSession(cookie.Value)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, err)
				return
			}
			userID = session.UserID
			// related to session an cookies updates:
			expTime, err := h.service.UserServiceInterface.ExtendSessionTimeout(cookie.Value)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("Cookie cannot be extended"))
				return
			}
			err = helpers.SessionCookieExtend(r, w, expTime)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, err)
				return
			}

		}

		field := getFiltersFieldFromURL(r.URL.Path)
		posts, err := h.service.PostServiceInterface.Filter(field, userID)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}

		for _, post := range posts {
			// getting the username for posts
			user, err := h.service.UserServiceInterface.GetUserByUserID(post.UserID)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, err)
				return
			}
			post.Username = user.Username

			// changing the format of the time
			post.CreatedTimeString = post.CreatedTime.Format("Jan 2, 2006 at 15:04")

			// assigning categories to each post
			temp_categories, err := h.service.PostServiceInterface.GetCategories(post.PostID)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, err)
				return
			}
			post.Categories = append(post.Categories, temp_categories...)
		}
		indexPath := "internal/web/templates/index.html"

		data := templateData{
			LoggedIn:      h.service.IsUserLoggedIn(r),
			AllPosts:      posts,
			AllCategories: []string{"Movie", "Game", "Book", "Others"}, // Initialize AllCategories with values
		}
		helpers.RenderTemplate(w, indexPath, data)
	default:
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Error in Post Reaction Handler"))
		return
	}
}

func getFiltersFieldFromURL(url string) string {
	return strings.Title(strings.TrimPrefix(url, "/filter/"))
}
