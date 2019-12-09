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

func (c *Client) Set(ctx context.Context) error {
	ctx, span := startSpan(ctx, "set")
	defer span.End()

	v, err := c.conn.Do("SET", "mykey", "いえーい", "nx")
	if err != nil {
		return err
	}
	fmt.Printf("set res: %+v", v)
	return nil
}
