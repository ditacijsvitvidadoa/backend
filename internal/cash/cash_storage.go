package cash

import (
	"encoding/json"
	"fmt"
	"github.com/ditacijsvitvidadoa/backend/internal/entities"
	"github.com/gomodule/redigo/redis"
	"log"
	"strings"
)

func SaveSessionToRedis(redisClient redis.Conn, cookieSession, token string) error {

	_, err := redisClient.Do("SET", cookieSession, token)
	if err != nil {
		return fmt.Errorf("could not save session to Redis: %v", err)
	}

	_, err = redisClient.Do("EXPIRE", cookieSession, 72*3600)
	if err != nil {
		return fmt.Errorf("could not set expiration for session in Redis: %v", err)
	}

	fmt.Println("Session saved in Redis for session:", cookieSession)
	return nil
}

func GetSessionFromRedis(redisClient redis.Conn, sessionID string) (string, error) {
	result, err := redis.String(redisClient.Do("GET", sessionID))
	if err == redis.ErrNil {
		return "", fmt.Errorf("session not found")
	} else if err != nil {
		return "", fmt.Errorf("redis error: %s", err)
	}

	return result, nil
}

func GetAllSessionsFromRedis(redisClient redis.Conn, userID string) ([]string, error) {
	pattern := fmt.Sprintf("session:%s:*", userID)

	keys, err := redis.Strings(redisClient.Do("KEYS", pattern))
	if err != nil {
		return nil, fmt.Errorf("Error fetching session keys: %s", err)
	}

	return keys, nil
}

func DeleteSessionByToken(redisClient redis.Conn, token string) error {
	keys, err := redis.Strings(redisClient.Do("KEYS", "*"))
	if err != nil {
		return fmt.Errorf("could not retrieve session keys: %v", err)
	}

	for _, sessionKey := range keys {
		sessionToken, err := redis.String(redisClient.Do("GET", sessionKey))
		if err == redis.ErrNil {
			continue
		} else if err != nil {
			return fmt.Errorf("redis error while fetching session: %s", err)
		}

		if sessionToken == token {
			_, err := redisClient.Do("DEL", sessionKey)
			if err != nil {
				return fmt.Errorf("could not delete session: %v", err)
			}
			fmt.Println("Deleted session for token:", token)
			return nil
		}
	}

	return fmt.Errorf("session with token %s not found", token)
}

func GetCitiesFromRedis(redisClient redis.Conn, searchQuery string) ([]entities.City, error) {
	var cities []entities.City
	cursor := "0"
	searchQuery = strings.ToLower(searchQuery)

	for {
		// Сканируем ключи с префиксом "city."
		reply, err := redis.Values(redisClient.Do("SCAN", cursor, "MATCH", "city.*", "COUNT", 100))
		if err != nil {
			return nil, fmt.Errorf("could not retrieve cities: %v", err)
		}

		var newCursor string
		var keys []string
		if _, err := redis.Scan(reply, &newCursor, &keys); err != nil {
			return nil, fmt.Errorf("could not scan keys: %v", err)
		}

		for _, key := range keys {
			cityData, err := redis.Bytes(redisClient.Do("GET", key))
			if err != nil {
				log.Printf("Error getting data for key %s: %v", key, err)
				continue
			}

			var city entities.City
			if err := json.Unmarshal(cityData, &city); err != nil {
				log.Printf("Error unmarshalling data for key %s: %v", key, err)
				continue
			}

			if strings.Contains(strings.ToLower(city.Description), searchQuery) {
				cities = append(cities, city)

				if len(cities) >= 20 {
					return cities, nil
				}
			}
		}

		if newCursor == "0" {
			break
		}
		cursor = newCursor
	}

	return cities, nil
}
