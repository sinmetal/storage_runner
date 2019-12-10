package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Client struct {
	pool *redis.Pool
}

func (c *Client) Close() error {
	return c.pool.Close()
}

func (c *Client) GetConn() redis.Conn {
	return c.pool.Get()
}

func NewClient(address string) (*Client, error) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
	return &Client{pool: pool}, nil
}

func Set(ctx context.Context, conn redis.Conn, key string, value string) error {
	ctx, span := startSpan(ctx, "set")
	defer span.End()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	return nil
}

func Get(ctx context.Context, conn redis.Conn, key string) (string, error) {
	ctx, span := startSpan(ctx, "get")
	defer span.End()

	v, err := conn.Do("GET", key)
	s, err := redis.String(v, err)
	if err != nil {
		if err == redis.ErrNil {
			// noop
		} else {
			return "", err
		}
	}

	return s, nil
}
