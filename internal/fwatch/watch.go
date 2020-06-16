package fwatch

import (
	"log"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func watchTargets(targets []Target) []*chan bool {
	dones := make([]*chan bool, len(targets))
	for i, t := range targets {
		done := make(chan bool)
		script := t.Script
		handler := func() bool {
			cmd := exec.Command("sh", "-c", script)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if IsService {
				cmd.Stdout = PidFile
				cmd.Stderr = PidFile
			}
			err := cmd.Run()
			if err != nil {
				print("exec error: %v", err)
			}
			return true
		}
		go watch(t.Path, t.Type, handler, done)
		dones[i] = &done
	}
	return dones
}

func watch(path string, types []string, handler func() bool, done chan bool) {
	print("start watch: path: %s", path)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if matchEvent(event, types) {
					if !handler() {
						break
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				print("exec error: %f", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		print("error: %v", err)
		return
	}
	<-done
	print("stop watch: path: %s", path)
}

func matchEvent(event fsnotify.Event, types []string) bool {
	switch true {
	case event.Op&fsnotify.Create == fsnotify.Create:
		return contains(types, Create)
	case event.Op&fsnotify.Remove == fsnotify.Remove:
		return contains(types, Remove)
	case event.Op&fsnotify.Rename == fsnotify.Rename:
		return contains(types, Rename)
	case event.Op&fsnotify.Write == fsnotify.Write:
		return contains(types, Write)
	case event.Op&fsnotify.Chmod == fsnotify.Chmod:
		return contains(types, Chmod)
	}
	return false
}

func contains(arr []string, s EventType) bool {
	for _, v := range arr {
		if s.Eq(v) {
			return true
		}
	}
	return false
}
