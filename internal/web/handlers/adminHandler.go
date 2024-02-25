package handlers

import (
	"errors"
	"forum/internal/models"
	"forum/internal/web/handlers/helpers"
	"net/http"
)

func (h *Handler) AdminMainPageHandler(w http.ResponseWriter, r *http.Request) {
	adminLoginPath := "internal/web/templates/adminPage.html"
	type templateData struct {
		AllRequests []*models.User
	}

	switch r.Method {
	case "GET":
		pendingUsers, err := h.service.UserServiceInterface.GetUsersByRole("pending")
		if err != nil {
			// fmt.Println("WHEN NO PENDING??")
			helpers.ErrorHandler(w, http.StatusBadRequest, errors.New("PEDNING USERS were not found"))
		}
		// fmt.Println("PENDING: ", pendingUsers, "Err: ", err)
		data := templateData{
			AllRequests: pendingUsers,
		}
		helpers.RenderTemplate(w, adminLoginPath, data)
		return
	// case "POST":
	// 	r.ParseForm()
	// 	fmt.Println("Raw Form Data:", r.Form)
	// 	fmt.Println("USER ID: ", r.FormValue("userId"), "ACTUION: ", r.FormValue("action"))
	default:
		helpers.ErrorHandler(w, http.StatusMethodNotAllowed, errors.New("in Admin Page Handler"))
		return
	}
}
