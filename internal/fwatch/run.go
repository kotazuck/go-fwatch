package fwatch

import (
	"os"
	"os/signal"
	"time"
)

// Run - run
func Run(path string) error {
	config, err := Load(path)
	if err != nil {
		return err
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	dones := watchTargets(config.Targets)

	doneConf := make(chan bool)
	go watch(config.path, []string{"write"}, func() bool {
		print("change config file: %s", config.path)
		for _, done := range dones {
			*done <- true
		}
		config, err := Load(path)
		if err != nil {
			print("error: load config file: %v", err)
			return true
		}
		dones = watchTargets(config.Targets)
		return true
	}, doneConf)

	<-quit

	print("begin remove all handler")

	go func() {
		doneConf <- true
		for _, done := range dones {
			*done <- true
		}
	}()
	time.Sleep(3 * time.Second)

	print("finish remove all handler")

	if IsService {
	}
	return nil
}
