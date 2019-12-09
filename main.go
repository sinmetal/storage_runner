package main

import (
	"fmt"

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

	endCh := make(chan error, 10)

	goSetRedis(rc, 3, endCh)
	goGetRedis(rc, 3, endCh)

	err = <-endCh
	fmt.Printf("BOMB %+v\n", err)
}
