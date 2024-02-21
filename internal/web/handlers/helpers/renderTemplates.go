package helpers

import (
	"errors"
	"html/template"
	"log"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, htmlTemplatePath string, resp interface{}) {
	temp, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, errors.New("problem parsing template"))
		log.Println(err)
		return
	}

	err = temp.Execute(w, resp)

	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, errors.New("problem executing template"))
		log.Println(err)
		return
	}
}
