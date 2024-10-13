package cash

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"os"
)

type RedisClient struct {
	Conn redis.Conn
}

func RedisConnection() (*RedisClient, error) {
	host, hostIsOk := os.LookupEnv("REDISDB_HOST")
	if !hostIsOk {
		host = "localhost"
	}

	port, portIsOk := os.LookupEnv("REDISDB_PORT")
	if !portIsOk {
		port = "6379"
	}

	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to Redis at", address)

	return &RedisClient{Conn: conn}, nil
}

func (rc *RedisClient) Close() error {
	if rc.Conn != nil {
		err := rc.Conn.Close()
		if err != nil {
			return err
		}
		fmt.Println("Connection to Redis closed")
	}
	return nil
}
