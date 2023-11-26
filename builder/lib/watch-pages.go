package lib

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path"
	"path/filepath"
)

type PageModificationEventType int

const (
	PageAdded PageModificationEventType = iota + 1
	PageDeleted
)

type PageWatchingEvent struct {
	Event           PageModificationEventType
	PathFromPageDir string
}

type PageWatcher struct {
	watcher *fsnotify.Watcher
	Events  chan PageWatchingEvent

	watchingPath string
}

func (w *PageWatcher) Init() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	w.watcher = watcher
	w.Events = make(chan PageWatchingEvent, 10)
}

func (w *PageWatcher) SetPageDir(packageDir, pageDir string) {
	w.watchingPath = path.Join(packageDir, pageDir)
	w.watcher.Add(w.watchingPath)
}

func (w *PageWatcher) Run() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)

			targetFilePageName, _ := filepath.Rel(w.watchingPath, event.Name)
			pageFileExt := filepath.Ext(targetFilePageName)

			// Page 는 tsx 만
			if pageFileExt != ".tsx" {
				return
			}

			if event.Has(fsnotify.Create) {
				log.Println("Create file:", targetFilePageName)
				w.Events <- PageWatchingEvent{
					Event:           PageAdded,
					PathFromPageDir: targetFilePageName,
				}
			}
			if event.Has(fsnotify.Remove) {
				log.Println("Remove file:", targetFilePageName)
				w.Events <- PageWatchingEvent{
					Event:           PageDeleted,
					PathFromPageDir: targetFilePageName,
				}
			}
			if event.Has(fsnotify.Rename) {
				log.Println("Rename file:", targetFilePageName)
				w.Events <- PageWatchingEvent{
					Event:           PageDeleted,
					PathFromPageDir: targetFilePageName,
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (w *PageWatcher) Close() {
	w.watcher.Close()
}
