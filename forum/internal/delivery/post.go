package delivery

import (
	"errors"
	"fmt"
	m "forum/internal/models"
	s "forum/internal/service"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func IDFromURL(url, prefix string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(url, prefix))
}

func (h *Handler) posts(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(m.User)
	postID, err := IDFromURL(r.URL.Path, "/posts/")
	if err != nil {
		h.errorPage(w, http.StatusNotFound, fmt.Errorf("error getting post ID: %s", err))
		return
	}

	switch r.Method {
	case http.MethodGet:
		post, err := h.services.Post.PostById(postID, user.ID)
		if err != nil {
			if errors.Is(err, s.ErrNoPost) {
				h.errorPage(w, http.StatusNotFound, nil)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		comments, err := h.services.Commentary.CommentsByPostID(postID, user.ID)
		if err != nil {
			log.Printf("error getting comments by post ID: %s", err)
		}

		data := m.TemplateData{
			User:     user,
			Post:     post,
			Comments: comments,
		}

		if err := h.tmpl.ExecuteTemplate(w, "post-comments", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
		}
	case http.MethodPost:
		if user == (m.User{}) {
			h.errorPage(w, http.StatusUnauthorized, nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		commentContent, ok := r.Form["comment"]
		if !ok {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		comment := m.Comment{
			UserID:  user.ID,
			PostID:  postID,
			Content: commentContent[0],
		}

		if err := h.services.Commentary.CreateComment(comment); err != nil {
			if errors.Is(err, s.ErrEmptyComment) {
				h.errorPage(w, http.StatusBadRequest, err)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(m.User)
	if user == (m.User{}) {
		h.errorPage(w, http.StatusUnauthorized, nil)
		return
	}

	switch r.Method {
	case http.MethodGet:
		data := m.TemplateData{
			User: user,
		}

		if err := h.tmpl.ExecuteTemplate(w, "create-post", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
		title, ok1 := r.Form["title"]
		content, ok2 := r.Form["content"]
		category, ok3 := r.Form["category"]

		if !ok1 || !ok2 || !ok3 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		post := m.Post{
			Title:      title[0],
			AuthorID:   user.ID,
			Content:    content[0],
			Categories: category,
		}

		if err := h.services.Post.CreatePost(post); err != nil {
			if errors.Is(err, s.ErrEmptyPost) {
				h.errorPage(w, http.StatusBadRequest, err)
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) reaction(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(contextKeyUser).(m.User)

	if user == (m.User{}) {
		h.errorPage(w, http.StatusUnauthorized, nil)
		return
	}

	if r.Method == http.MethodGet {
		h.errorPage(w, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodPost {
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	reaction, ok := r.Form["react"]
	if !ok {
		h.errorPage(w, http.StatusBadRequest, nil)
		return
	}

	id, err := IDFromURL(r.URL.Path, "/posts/react/")
	if err != nil {
		h.errorPage(w, http.StatusNotFound, fmt.Errorf("error getting post ID: %s", err))
		return
	}

	if err := h.services.Reaction.ReactToPost(id, user.ID, reaction[0]); err != nil {
		if err.Error() == http.StatusText(http.StatusBadRequest) {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}
		h.errorPage(w, http.StatusInternalServerError, nil)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/posts/%v", id), http.StatusSeeOther)
}
