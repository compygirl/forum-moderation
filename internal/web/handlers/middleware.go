package handlers

import (
	"errors"
	"forum/internal/web/handlers/helpers"
	"net/http"
)

const cookieName = "session_id"

func (h *Handler) CheckCookieMiddleware(someHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil && c != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Error with Cookie11!!!"))
			return
		}
		if c != nil {
			if !h.service.IsTokenExist(c.Value) {
				helpers.SessionCookieExpire(w)
			}
		}
		someHandler.ServeHTTP(w, r)
	})
}

func (h *Handler) OnlyUnauthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil && c != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Error with Cookie"))
			return
		}
		if c != nil {
			http.Redirect(w, r, "/", 302)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Handler) NeedAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil && c != nil {
			helpers.ErrorHandler(w, http.StatusInternalServerError, errors.New("Error with Cookie"))
			return
		}
		if c == nil {
			http.Redirect(w, r, "/login", 302)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
