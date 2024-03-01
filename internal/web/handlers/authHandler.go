package handlers

import (
	"errors"
	"fmt"
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
		var userRole string
		admin := r.FormValue("admin") == "on"
		if admin {
			fmt.Println("Registration of ADMIN")
			userRole = "admin"
		} else {
			userRole = "user"
		}

		user := &models.User{
			FirstName:  r.FormValue("firstName"),
			SecondName: r.FormValue("secondName"),
			Username:   r.FormValue("username"),
			Email:      r.FormValue("email"),
			Password:   string(psw),
			Role:       userRole,
		}

		statusCode, id, err := h.service.UserServiceInterface.CreateUser(user)
		if err != nil {
			helpers.ErrorHandler(w, statusCode, err)
			return
		}
		user.UserUserID = id
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
		admin := r.FormValue("admin") == "on"
		// check users credentials and handle the admin level
		// if admin {
		// } else {
		// fmt.Println("ADMING BEFORE: ", admin)
		session, err := h.service.UserServiceInterface.Login(email, password, admin)
		if err != nil {
			helpers.ErrorHandler(w, http.StatusForbidden, err)
			return
		} // else {

		helpers.SessionCookieSet(w, session.Token, session.ExpTime)
		// fmt.Println("ADMING AFTER: ", admin)
		if admin {
			// fmt.Println("ADMIN LOGINNING")
			http.Redirect(w, r, "/admin_page", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		//}
		// }

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
