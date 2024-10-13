package cookie

import (
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/cash"
	"github.com/ditacijsvitvidadoa/backend/internal/utils"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CheckCookie(cashConn redis.Conn, r *http.Request) (bool, string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			return false, "", nil
		}
		return false, "", err
	}

	decodedValue, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return false, "", fmt.Errorf("Error decoding cookie value: %v", err)
	}

	parts := strings.Split(decodedValue, " ")
	if len(parts) != 2 {
		return false, "", fmt.Errorf("Invalid cookie format")
	}

	userIdPart := parts[0]
	tokenPart := parts[1]

	userIDParts := strings.Split(userIdPart, "=")
	if len(userIDParts) != 2 {
		return false, "", fmt.Errorf("Invalid userId format")
	}
	userID := userIDParts[1]

	tokenParts := strings.Split(tokenPart, "=")
	if len(tokenParts) != 2 {
		return false, "", fmt.Errorf("Invalid token format")
	}
	token := tokenParts[1]

	sessions, err := cash.GetAllSessionsFromRedis(cashConn, userID)
	if err != nil {
		return false, "", err
	}

	for _, sessionKey := range sessions {
		storedToken, err := cash.GetSessionFromRedis(cashConn, sessionKey)
		if err != nil {
			return false, "", fmt.Errorf("Error getting token from Redis: %s", err)
		}

		if storedToken == token {
			claims, err := utils.ValidateJWT(token)
			if err != nil {
				return false, "", fmt.Errorf("Invalid token: %v", err)
			}

			if claims.ExpiresAt < time.Now().Unix() {
				return false, "", fmt.Errorf("Token expired")
			}

			return true, userID, nil
		}
	}

	return false, "", fmt.Errorf("Invalid session")
}
