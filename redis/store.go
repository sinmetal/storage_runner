package redis

import (
	"context"
	"fmt"
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

	v, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	fmt.Printf("set res: key:%s:%+v\n", key, v)
	return nil
}

func Get(ctx context.Context, conn redis.Conn, key string) error {
	ctx, span := startSpan(ctx, "get")
	defer span.End()

	v, err := conn.Do("GET", key)
	if err != nil {
		return err
	}
	fmt.Printf("get res: key:%s:%+v\n", key, v)
	return nil
}
