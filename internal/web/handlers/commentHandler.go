package handlers

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/web/handlers/helpers"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) DisplayCommentsHandler(w http.ResponseWriter, r *http.Request) {
	commentsPath := "internal/web/templates/comments.html"
	type templateData struct {
		LoggedIn    bool
		ThePost     *models.Post
		AllComments []*models.Comment
	}

	switch r.Method {
	case "GET":
		postId := getPostIDFromURL(r.URL.Path)
		comments, err := h.service.CommentServiceInterface.GetAlCommentsForPost(postId)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}

		for _, comment := range comments {
			user, err := h.service.UserServiceInterface.GetUserByUserID(comment.UserID)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, err)
				return
			}
			comment.Username = user.Username
			comment.CreatedTimeString = comment.CreatedTime.Format("Jan 2, 2006 at 15:04")
		}

		post, err := h.service.PostServiceInterface.GetPostByID(postId)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}

		// get username
		user, err := h.service.UserServiceInterface.GetUserByUserID(post.UserID)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}
		post.Username = user.Username

		// change the time format
		post.CreatedTimeString = post.CreatedTime.Format("Jan 2, 2006 at 15:04")

		cookie := helpers.SessionCookieGet(r)
		if cookie != nil {
			expTime, err := h.service.UserServiceInterface.ExtendSessionTimeout(cookie.Value)
			if err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("The Time cannot be extended"))
				return
			}
			if err := helpers.SessionCookieExtend(r, w, expTime); err != nil {
				helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Cookie cannot be extended"))
				return
			}
		}

		helpers.RenderTemplate(w, commentsPath, templateData{h.service.IsUserLoggedIn(r), post, comments})
	default:
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Error in DisplayCommentsHandler"))
		return
	}
}

func getPostIDFromURL(url string) int {
	strID := strings.TrimPrefix(url, "/comments/")
	id, err := strconv.Atoi(strID)
	if err != nil {
		return -1
	}
	return id
}

func (h *Handler) CreateCommentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		cookie := helpers.SessionCookieGet(r)
		if cookie == nil {
			helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Unauthenticated user during creating comment"))
			return
		}

		session, err := h.service.UserServiceInterface.GetSession(cookie.Value)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Session doesn't exist"))
			return
		}

		postId, err := strconv.Atoi(r.FormValue("postId"))
		if err != nil {
			helpers.ErrorHandler(w, http.StatusBadRequest, errors.New("Converstion of PostID is not allowed"))
			return
		}

		comment := &models.Comment{
			UserID:  session.UserID,
			PostID:  postId,
			Content: r.FormValue("commentcontent"),
		}

		statusCode, _, err := h.service.CommentServiceInterface.CreateComment(comment)
		if err != nil {
			helpers.ErrorHandler(w, statusCode, err)
			return
		}

		expTime, err := h.service.UserServiceInterface.ExtendSessionTimeout(cookie.Value)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("The Time cannot be extended"))
			return
		}
		if err := helpers.SessionCookieExtend(r, w, expTime); err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Cookie cannot be extended"))
			return
		}
		http.Redirect(w, r, "/comments/"+fmt.Sprint(comment.PostID), http.StatusSeeOther)
		return
	default:
		helpers.ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("Error in Comment creation Handler"))
		return

	}
}

func (h *Handler) ReactOnCommentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil {
			helpers.ErrorHandler(w, http.StatusBadRequest, errors.New("Coversion of commentID failed"))
			return
		}

		currReaction, err := strconv.Atoi(r.FormValue("type"))
		if err != nil {
			helpers.ErrorHandler(w, http.StatusBadRequest, errors.New("Conversion of reaction type failed"))
			return
		}

		cookie := helpers.SessionCookieGet(r)
		if cookie == nil {
			helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Cookie cannot be reseived in Comment Reaction Handler"))
			return
		}

		session, err := h.service.UserServiceInterface.GetSession(cookie.Value)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, fmt.Errorf("Get Session in Reactin on Comment: %w", err))
			return
		}

		postId, err := strconv.Atoi(r.FormValue("postId"))
		if err != nil {
			helpers.ErrorHandler(w, http.StatusBadRequest, errors.New("Converstion of PostID is not allowed"))
			return
		}
		if err := h.service.CommentServiceInterface.UpdateReaction(currReaction, commentID, session.UserID); err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		}

		// related to session an cookies updates:
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
		http.Redirect(w, r, "/comments/"+fmt.Sprint(postId), http.StatusSeeOther)
		return
	default:
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Error in Comment Reaction Handler"))
		return
	}
}
