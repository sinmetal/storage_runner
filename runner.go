package main

import (
	"context"
	"sync"
	"time"

	"github.com/sinmetal/storage_runner/redis"
)

func goSetRedis(rc *redis.Client, goroutine int, endCh chan<- error) {
	go func() {
		for {
			var wg sync.WaitGroup
			for i := 0; i < goroutine; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					ctx := context.Background()

					id := NewNarrowRandom()

					ctx, span := startSpan(ctx, "setRedis")
					defer span.End()

					var cancel context.CancelFunc
					if _, hasDeadline := ctx.Deadline(); !hasDeadline {
						ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
						defer cancel()
					}

					if err := rc.Set(ctx, id, id); err != nil {
						endCh <- err
					}
				}(i)
			}
			wg.Wait()
		}
	}()
}

func goGetRedis(rc *redis.Client, goroutine int, endCh chan<- error) {
	go func() {
		for {
			var wg sync.WaitGroup
			for i := 0; i < goroutine; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					ctx := context.Background()
					id := NewNarrowRandom()

					ctx, span := startSpan(ctx, "getRedis")
					defer span.End()

					var cancel context.CancelFunc
					if _, hasDeadline := ctx.Deadline(); !hasDeadline {
						ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
						defer cancel()
					}

					if err := rc.Get(ctx, id); err != nil {
						endCh <- err
					}
				}(i)
			}
			wg.Wait()
		}
	}()
}
