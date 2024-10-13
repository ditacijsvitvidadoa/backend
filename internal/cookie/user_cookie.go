package cookie

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func SetCookie(w http.ResponseWriter, session string, token string) {
	expirationTime := time.Now().Add(72 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:     session,
		Value:    token,
		Path:     "/",
		Expires:  expirationTime,
		MaxAge:   72 * 3600,
		HttpOnly: true,
		Secure:   false,
	})
}

func GetUserIDFromCookie(cookieValue string) (string, error) {
	parts := strings.Split(cookieValue, " ")
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
		Name:     "session",       // Название куки
		Value:    "",              // Устанавливаем пустое значение
		Path:     "/",             // Путь, к которому относится кука
		Expires:  time.Unix(0, 0), // Устанавливаем время истечения в прошлом
		MaxAge:   -1,              // Указываем, что кука должна быть удалена
		HttpOnly: true,            // Защита от доступа через JavaScript
		Secure:   true,            // Устанавливайте true, если используете HTTPS
	})
}
