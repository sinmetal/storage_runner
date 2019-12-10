package redis

import (
	"context"
	"testing"
)

func TestClient_Set(t *testing.T) {
	c, err := NewClient("127.0.0.1:6379")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Set(context.Background(), "hoge", "hoge"); err != nil {
		t.Fatal(err)
	}

}

func TestClient_Get(t *testing.T) {
	c, err := NewClient("127.0.0.1:6379")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := c.Set(ctx, "hoge", "hoge"); err != nil {
		t.Fatal(err)
	}

	if err := c.Get(ctx, "hoge"); err != nil {
		t.Fatal(err)
	}
}
