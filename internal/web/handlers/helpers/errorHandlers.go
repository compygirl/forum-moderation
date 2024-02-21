package helpers

import (
	models "forum/internal/models"
	"html/template"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, errorNum int, errDetails error) {
	var resp models.ErrorResponse
	resp.ErrorNum = errorNum
	resp.ErrorMessage = http.StatusText(errorNum) + "\n" + errDetails.Error()
	w.WriteHeader(errorNum)

	temp, err := template.ParseFiles("./internal/web/templates/errors.html")

	err = temp.Execute(w, resp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
