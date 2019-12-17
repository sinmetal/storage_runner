package main

import (
	"context"
	"fmt"
	"github.com/morikuni/failure"
	"github.com/sinmetal/storage_runner/stats"
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
						ctx, cancel = context.WithTimeout(ctx, 500*time.Millisecond)
						defer cancel()
					}

					retCh := make(chan error)
					go func() {
						ret := redis.Set(ctx, conn, id, id)
						select {
						case <-ctx.Done():
						case retCh <- ret:
						}
					}()
					select {
					case <-ctx.Done():
						if err := stats.CountRedisStatus(ctx, "SET TIMEOUT"); err != nil {
							endCh <- err
						}
					case err := <-retCh:
						serr := stats.CountRedisStatus(ctx, "SET NG")
						if serr != nil {
							err = failure.Wrap(err, failure.Messagef("failed stats. err=%+v", serr))
						}
						if err != nil {
							endCh <- err
						}
					}

					if err := stats.CountRedisStatus(ctx, "SET OK"); err != nil {
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
						ctx, cancel = context.WithTimeout(ctx, 500*time.Millisecond)
						defer cancel()
					}

					retCh := make(chan error)
					go func() {
						_, ret := redis.Get(ctx, conn, id)
						select {
						case <-ctx.Done():
						case retCh <- ret:
						}
					}()
					select {
					case <-ctx.Done():
						if err := stats.CountRedisStatus(ctx, "GET TIMEOUT"); err != nil {
							endCh <- err
						}
					case err := <-retCh:
						serr := stats.CountRedisStatus(ctx, "GET NG")
						if serr != nil {
							err = failure.Wrap(err, failure.Messagef("failed stats. err=%+v", serr))
						}
						if err != nil {
							endCh <- err
						}
					}

					if err := stats.CountRedisStatus(ctx, "GET OK"); err != nil {
						endCh <- err
					}
				}(i)
			}
			wg.Wait()
		}
	}()
}
