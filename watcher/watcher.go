package watcher

import (
    "sync"
    "time"

    "github.com/fsnotify/fsnotify"
)

type Handler struct {
    OnWrite  func(string)
    OnDelete func(string)
}

func Start(path string, h Handler) {
    w, _ := fsnotify.NewWatcher()

    debounce := map[string]*time.Timer{}
    var mu sync.Mutex

    go func() {
        for {
            select {
            case e := <-w.Events:

                if e.Op&(fsnotify.Write|fsnotify.Create) != 0 {
                    mu.Lock()

                    if t, ok := debounce[e.Name]; ok {
                        t.Stop()
                    }

                    debounce[e.Name] = time.AfterFunc(500*time.Millisecond, func() {
                        h.OnWrite(e.Name)

                        mu.Lock()
                        delete(debounce, e.Name)
                        mu.Unlock()
                    })

                    mu.Unlock()
                }

                if e.Op&fsnotify.Remove != 0 {
                    go h.OnDelete(e.Name)
                }

            case <-w.Errors:
            }
        }
    }()

    w.Add(path)
}