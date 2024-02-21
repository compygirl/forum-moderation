package helpers

import (
	"fmt"
	"net/http"
	"time"
)

const cookieName = "session_id"

func SessionCookieGet(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}
	return cookie
}

func SessionCookieExpire(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, &cookie)
}

func SessionCookieSet(w http.ResponseWriter, token string, expirationTime time.Time) {
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   token,
		Path:    "/",
		Expires: expirationTime,
	}
	http.SetCookie(w, &cookie)
}

// TODO: Handle Errors where we are using this fucntion
func SessionCookieExtend(r *http.Request, w http.ResponseWriter, expirationTime time.Time) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return fmt.Errorf("SessionCookieExtend: %w", err)
	}
	cookie.Expires = expirationTime
	cookie.Path = "/"
	http.SetCookie(w, cookie)
	return nil
}
