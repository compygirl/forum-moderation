package handlers

import (
	"errors"
	"forum/internal/models"
	helpers "forum/internal/web/handlers/helpers"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	registerPath := "internal/web/templates/registration.html"

	switch r.Method {
	case "GET":
		helpers.RenderTemplate(w, registerPath, nil)
		return
	case "POST":
		psw, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
		user := &models.User{
			FirstName:  r.FormValue("firstName"),
			SecondName: r.FormValue("secondName"),
			Username:   r.FormValue("username"),
			Email:      r.FormValue("email"),
			Password:   string(psw),
		}

		statusCode, id, err := h.service.UserServiceInterface.CreateUser(user)
		if err != nil {
			helpers.ErrorHandler(w, statusCode, err)
			return
		}
		user.UserID = id
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	default:
		helpers.ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("Error in Registration Handler"))
		return
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginPath := "internal/web/templates/login.html"

	switch r.Method {
	case "GET":
		helpers.RenderTemplate(w, loginPath, nil)
		return
	case "POST":

		email := r.FormValue("email")
		password := r.FormValue("password")

		session, err := h.service.UserServiceInterface.Login(email, password)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusBadRequest, err)
			return
		} else {
			helpers.SessionCookieSet(w, session.Token, session.ExpTime)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	default:
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Error in Login Handler"))
		return
	}
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cookie := helpers.SessionCookieGet(r)
		if cookie == nil {
			helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Conversion of postID failed"))
			return
		}

		//??
		if err := h.service.UserServiceInterface.Logout(cookie.Value); err != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, err)
			return
		} else {
			helpers.SessionCookieExpire(w)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	default:
		helpers.ErrorHandler(w, http.StatusUnauthorized, errors.New("Error in Logout Handler"))
		return
	}
}
