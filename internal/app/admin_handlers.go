package app

import (
	"github.com/ditacijsvitvidadoa/backend/internal/storage/requests"
	"net/http"
	"time"
)

func (a *App) loginAdminPanel(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("Email")
	password := r.FormValue("Password")
	phone := r.FormValue("Phone")

	ok, err := requests.ValidateAdminCredentials(a.client, email, password, phone)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		sendError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	cookie := &http.Cookie{
		Name:     "admin-session",
		Value:    "true",
		Expires:  expirationTime,
		MaxAge:   24 * 3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, cookie)

	sendOk(w)
}

func (a *App) checkAdminSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("admin-session")
	if err != nil {
		if err == http.ErrNoCookie {
			sendError(w, http.StatusUnauthorized, "no admin session cookie found")
		} else {
			sendError(w, http.StatusInternalServerError, "failed to read cookie")
		}
		return
	}

	if cookie.Value != "true" {
		sendError(w, http.StatusUnauthorized, "invalid admin session")
		return
	}

	sendOk(w)
}
