package main

import (
	"context"
	"fmt"
	"github.com/sinmetal/storage_runner/metrics"
	"sync"
	"time"

	"github.com/sinmetal/storage_runner/redis"
)

func GoSetRedis(rc *redis.Client, goroutine int, endCh chan<- error) {
	go func() {
		for {
			var wg sync.WaitGroup
			for i := 0; i < goroutine; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					conn := rc.GetConn()
					defer func() {
						if err := conn.Close(); err != nil {
							fmt.Printf("failed redis.Conn.Close().err=%+v\n", err)
						}
					}()

					ctx := context.Background()

					id := NewNarrowRandom()

					ctx, span := startSpan(ctx, "setRedis")
					defer span.End()

					var cancel context.CancelFunc
					if _, hasDeadline := ctx.Deadline(); !hasDeadline {
						ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
						defer cancel()
					}

					if err := redis.Set(ctx, conn, id, id); err != nil {
						endCh <- err
					}
					if err := metrics.CountStatus(ctx, "SET OK"); err != nil {
						endCh <- err
					}
				}(i)
			}
			wg.Wait()
		}
	}()
}

func GoGetRedis(rc *redis.Client, goroutine int, endCh chan<- error) {
	go func() {
		for {
			var wg sync.WaitGroup
			for i := 0; i < goroutine; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					conn := rc.GetConn()
					defer func() {
						if err := conn.Close(); err != nil {
							fmt.Printf("failed redis.Conn.Close().err=%+v\n", err)
						}
					}()

					ctx := context.Background()
					id := NewNarrowRandom()

					ctx, span := startSpan(ctx, "getRedis")
					defer span.End()

					var cancel context.CancelFunc
					if _, hasDeadline := ctx.Deadline(); !hasDeadline {
						ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
						defer cancel()
					}

					_, err := redis.Get(ctx, conn, id)
					if err != nil {
						endCh <- err
					}

					if err := metrics.CountStatus(ctx, "GET OK"); err != nil {
						endCh <- err
					}
				}(i)
			}
			wg.Wait()
		}
	}()
}
