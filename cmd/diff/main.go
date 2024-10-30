package main

import (
	ga "github.com/sethvargo/go-githubactions"

	"tf-plan-diff/config"
	"tf-plan-diff/summary"
)

func run() error {
	action := ga.New()

	cfg, err := config.ConfigFromAction(action)
	if err != nil {
		return err
	}

	return summary.Run(cfg)
}

func main() {
	err := run()
	if err != nil {
		ga.Fatalf("%v", err)
	}
}
