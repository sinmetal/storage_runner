package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Client struct {
	conn redis.Conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func NewClient() (*Client, error) {
	c, err := redis.Dial("tcp", "10.146.99.131:6379")
	if err != nil {
		return nil, err
	}
	return &Client{conn: c}, nil
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	ctx, span := startSpan(ctx, "set")
	defer span.End()

	v, err := c.conn.Do("SET", key, value, "nx")
	if err != nil {
		return err
	}
	fmt.Printf("set res: key:%s:%+v\n", key, v)
	return nil
}

func (c *Client) Get(ctx context.Context, key string) error {
	ctx, span := startSpan(ctx, "get")
	defer span.End()

	v, err := c.conn.Do("GET", key)
	if err != nil {
		return err
	}
	fmt.Printf("get res: key:%s:%+v\n", key, v)
	return nil
}
