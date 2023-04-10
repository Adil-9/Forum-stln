package delivery

import (
	"forum/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) myPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(models.User)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		posts, err := h.services.Post.UsersPosts(user.ID)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		data := models.TemplateData{
			User:  user,
			Posts: posts,
		}

		if err := h.tmpl.ExecuteTemplate(w, "home", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		if user == (models.User{}) {
			h.errorPage(w, http.StatusUnauthorized, nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		postID, ok1 := r.Form["postID"]
		react, ok2 := r.Form["react"]

		if !ok1 || !ok2 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		postid, err := strconv.Atoi(postID[0])
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, nil)
			return
		}

		if err := h.services.ReactToPost(postid, user.ID, react[0]); err != nil {
			if err.Error() == http.StatusText(http.StatusBadRequest) {
				h.errorPage(w, http.StatusBadRequest, nil)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, nil)
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) likedPosts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(models.User)
	if user == (models.User{}) {
		h.errorPage(w, http.StatusUnauthorized, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		posts, err := h.services.Post.LikedPosts(user.ID)
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		data := models.TemplateData{
			User:  user,
			Posts: posts,
		}

		if err := h.tmpl.ExecuteTemplate(w, "home", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		if user == (models.User{}) {
			h.errorPage(w, http.StatusUnauthorized, nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		postIDVal, ok1 := r.Form["postID"]
		react, ok2 := r.Form["react"]

		if !ok1 || !ok2 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		postID, err := strconv.Atoi(postIDVal[0])
		if err != nil {
			h.errorPage(w, http.StatusInternalServerError, nil)
			return
		}

		if err := h.services.ReactToPost(postID, user.ID, react[0]); err != nil {
			if err.Error() == http.StatusText(http.StatusBadRequest) {
				h.errorPage(w, http.StatusBadRequest, nil)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, nil)
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}
