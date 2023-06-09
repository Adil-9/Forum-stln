package delivery

import (
	"errors"
	"fmt"
	m "forum/internal/models"
	s "forum/internal/service"
	"net/http"
	"time"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := m.TemplateData{}
		if err := h.tmpl.ExecuteTemplate(w, "reg_page", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		email, ok1 := r.Form["email"]
		username, ok2 := r.Form["username"]
		password, ok3 := r.Form["password"]
		confirm, ok4 := r.Form["confirm_password"]

		if !ok1 || !ok2 || !ok3 || !ok4 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		user := m.User{
			Email:           email[0],
			Username:        username[0],
			Password:        password[0],
			ConfirmPassword: confirm[0],
		}

		if err := h.services.Authorization.CreateUser(user); err != nil {
			if errors.Is(err, s.ErrInvalidEmail) || errors.Is(err, s.ErrInvalidPassword) ||
				errors.Is(err, s.ErrInvalidUsername) || errors.Is(err, s.ErrUsernameTaken) ||
				errors.Is(err, s.ErrEmailTaken) {
				data := m.TemplateData{Error: m.ErrorMsg{Msg: fmt.Sprint(err.Error())}}
				if err := h.tmpl.ExecuteTemplate(w, "reg_page", data); err != nil {
					h.errorPage(w, http.StatusInternalServerError, err)
					return
				}
				return
			} else {
				h.errorPage(w, http.StatusInternalServerError, err)
				return
			}
		}
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := m.TemplateData{}
		if err := h.tmpl.ExecuteTemplate(w, "sign-in", data); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		username, ok1 := r.Form["username"]
		password, ok2 := r.Form["password"]

		if !ok1 || !ok2 {
			h.errorPage(w, http.StatusBadRequest, nil)
			return
		}

		session, err := h.services.Authorization.SetSession(username[0], password[0])
		if err != nil {
			if errors.Is(err, s.ErrWrongPasswordOrUser) {
				data := m.TemplateData{Error: m.ErrorMsg{Msg: fmt.Sprint(err.Error())}}
				if err := h.tmpl.ExecuteTemplate(w, "sign-in", data); err != nil {
					h.errorPage(w, http.StatusInternalServerError, err)
					return
				}
				return
			}
			h.errorPage(w, http.StatusInternalServerError, err)
			return
		}

		cookie := &http.Cookie{
			Name:    "session_token",
			Value:   session.Token,
			Path:    "/",
			Expires: session.ExpirationDate,
		}
		http.SetCookie(w, cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		h.errorPage(w, http.StatusMethodNotAllowed, nil)
	}
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.errorPage(w, http.StatusUnauthorized, err)
			return
		}
		h.errorPage(w, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Authorization.DeleteSession(cookie.Value); err != nil {
		h.errorPage(w, http.StatusInternalServerError, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
