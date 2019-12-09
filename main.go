package main

import (
	"context"
	"math/rand"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/sinmetal/gcpmetadata"
	"github.com/sinmetal/storage_runner/redis"
	"go.opencensus.io/trace"
)

func main() {
	project, err := gcpmetadata.GetProjectID()
	if err != nil {
		panic(err)
	}

	{
		exporter, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: project,
		})
		if err != nil {
			panic(err)
		}
		trace.RegisterExporter(exporter)
		// trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	rc, err := redis.NewClient()
	if err != nil {
		panic(err)
	}
	for {
		ctx := context.Background()
		err = rc.Set(ctx)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
}
