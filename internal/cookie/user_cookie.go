package cookie

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func SetCookie(w http.ResponseWriter, session string, userId string, token string) {
	expirationTime := time.Now().Add(72 * time.Hour)
	value := fmt.Sprintf("userId=%s|token=%s", userId, token)
	cookie := &http.Cookie{
		Name:     session,
		Value:    value,
		Path:     "/",
		Expires:  expirationTime,
		MaxAge:   72 * 3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	log.Printf("Cookie details: %+v\n", cookie)
}

func GetUserIDFromCookie(cookieValue string) (string, error) {
	parts := strings.Split(cookieValue, "|")
	if len(parts) != 2 {
		return "", fmt.Errorf("Invalid cookie format")
	}

	userIdPart := parts[0]
	fmt.Println("userIdPart", userIdPart)

	userIDParts := strings.Split(userIdPart, "=")
	if len(userIDParts) != 2 || userIDParts[0] != "userId" {
		return "", fmt.Errorf("Invalid userId format")
	}

	return userIDParts[1], nil
}

func GetSessionValue(r *http.Request, cookieName string) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no cookie found")
		}
		return "", err
	}
	return cookie.Value, nil
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})
}
