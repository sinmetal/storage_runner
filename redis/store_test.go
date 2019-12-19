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
	if err := Set(context.Background(), c.GetConn(), "hoge", "hoge"); err != nil {
		t.Fatal(err)
	}

}

func TestClient_Get(t *testing.T) {
	c, err := NewClient("127.0.0.1:6379")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := Set(ctx, c.GetConn(), "hoge", "hoge"); err != nil {
		t.Fatal(err)
	}

	_, err = Get(ctx, c.GetConn(), "hoge")
	if err != nil {
		t.Fatal(err)
	}
}
