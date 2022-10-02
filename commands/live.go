package commands

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

var watcher *fsnotify.Watcher

func Live(addr string) error {
	// Build content
	Build()

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/", http.FileServer(http.Dir("./build")))

		server := http.Server{
			Addr:    addr,
			Handler: requestLog(mux),
		}

		fmt.Println("serving on", addr)
		log.Fatal(server.ListenAndServe())
	}()

	// Watch files changes
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk(".", watchDir); err != nil {
		log.Fatal("error:", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("modified:", event.Name)
				Build()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	<-done

	return nil
}

func watchDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		if strings.HasPrefix(path, "pages") {
			return watcher.Add(path)
		}
	}

	return nil
}
