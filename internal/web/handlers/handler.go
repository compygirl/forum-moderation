package handlers

import (
	service "forum/internal/service"
	"net/http"
	"time"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	handlerObj := Handler{service: service}
	return &handlerObj
}

func (handler *Handler) InitRouter() *http.ServeMux {
	mux := http.NewServeMux()

	images := http.FileServer(http.Dir("./data/assets/images"))
	mux.Handle("/images/", http.StripPrefix("/images/", images))
	mux.HandleFunc("/", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.GetMainPage)))
	mux.HandleFunc("/registration", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.OnlyUnauthMiddleware(handler.RegistrationHandler))))
	mux.HandleFunc("/login", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.OnlyUnauthMiddleware(handler.CheckCookieMiddleware(handler.LoginHandler))))
	mux.HandleFunc("/logout", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.LogoutHandler))))
	mux.HandleFunc("/submit-post", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.CreatePostHandler))))
	mux.HandleFunc("/post/react", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.ReactOnPostHandler))))
	mux.HandleFunc("/comments/", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.DisplayCommentsHandler)))
	mux.HandleFunc("/submit-comment", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.CreateCommentsHandler))))
	mux.HandleFunc("/comment/react", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.ReactOnCommentHandler))))
	mux.HandleFunc("/filter/", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.FilterHandler)))
	mux.HandleFunc("/auth/google/in", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.OnlyUnauthMiddleware(handler.GoogleAuthHandler))))
	mux.HandleFunc("/auth/google/callback", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.OnlyUnauthMiddleware(handler.GoogleCallback))))
	mux.HandleFunc("/auth/github/in", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.OnlyUnauthMiddleware(handler.GithubAuthHandler))))
	mux.HandleFunc("/auth/github/callback", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.OnlyUnauthMiddleware(handler.GithubCallback))))
	// mux.HandleFunc("/moderator", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.ModeratorRequestHandler))))
	mux.HandleFunc("/moderator", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.ModeratorRequestHandler))))
	// mux.HandleFunc("/moderator", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.AdminLoginHandler))))

	mux.HandleFunc("/admin_page", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.AdminMainPageHandler))))
	mux.HandleFunc("/approve", NewRateLimiter(10, time.Minute).LimitMiddleware(handler.CheckCookieMiddleware(handler.NeedAuthMiddleware(handler.ApproveModeratorHandler))))
	return mux
}
