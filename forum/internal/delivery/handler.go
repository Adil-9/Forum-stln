package delivery

import (
	"forum/internal/service"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl     *template.Template
	services *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		tmpl:     template.Must(template.ParseGlob("templates/*.html")),
		services: service,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.middleware(h.homePage))
	mux.HandleFunc("/sign-out", h.logout)
	mux.HandleFunc("/sign-up", h.signUp)
	mux.HandleFunc("/sign-in", h.signIn)
	mux.HandleFunc("/my-posts", h.middleware(h.myPosts))
	mux.HandleFunc("/liked-posts", h.middleware(h.likedPosts))
	mux.HandleFunc("/posts/", h.middleware(h.posts))
	mux.HandleFunc("/posts/create", h.middleware(h.createPost))
	mux.HandleFunc("/posts/react/", h.middleware(h.reaction))
	mux.HandleFunc("/comment/react/", h.middleware(h.reactToComment))

	mux.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("templates/"))))
	return mux
}
