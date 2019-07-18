package helper

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func FileWatcher(wpath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("file modified:", event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("file created:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(wpath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
