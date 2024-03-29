package main

import (
	"fmt"
	"log"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/kelseyhightower/envconfig"
	"github.com/sinmetal/gcpmetadata"
	sr "github.com/sinmetal/storage_runner/redis"
	"github.com/sinmetal/storage_runner/stats"
	"go.opencensus.io/trace"
)

type EnvConfig struct {
	RedisAddress string `default:"127.0.0.1:6379"`
}

func main() {
	var env EnvConfig
	if err := envconfig.Process("storage", &env); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("ENV_CONFIG %+v\n", env)

	project, err := gcpmetadata.GetProjectID()
	if err != nil {
		log.Printf("ProjectID not found")
	}

	if project != "" {
		exporter, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: project,
		})
		if err != nil {
			panic(err)
		}
		trace.RegisterExporter(exporter)
		// trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

		{
			exporter := stats.InitExporter(project)
			stats.InitOpenCensusStats(exporter)
		}
	}

	rc, err := sr.NewClient(env.RedisAddress)
	if err != nil {
		panic(err)
	}

	endCh := make(chan error, 10)

	GoSetRedis(rc, 3, endCh)
	GoGetRedis(rc, 3, endCh)

	err = <-endCh
	fmt.Printf("BOMB %+v\n", err)
}
