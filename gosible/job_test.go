package gosible

import (
	"testing"
	"time"
)

func TestAnsibleJobRun(t *testing.T) {
	p := make(map[string]interface{})
	p["pattern"] = "localhost"
	p["moduleName"] = "ping"
	p["forks"] = "5"
	p["limit"] = "localhost"

	job, err := NewJob(p)
	if err != nil {
		panic(err)
	}

	err = job.Run()
	if err != nil {
		panic(err)
	}
}

func TestAnsibleJobKill(t *testing.T) {
	p := make(map[string]interface{})
	p["pattern"] = "localhost"
	p["moduleName"] = "shell"
	p["moduleArgs"] = "'sleep 10'"
	p["forks"] = "5"
	p["limit"] = "localhost"

	job, err := NewJob(p)
	if err != nil {
		panic(err)
	}

	err = job.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	if err = job.Kill(); err != nil {
		panic(err)
	}
}
