package handler

import (
	"extool/base"
	"extool/scanner"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

type Walker struct {
	walkFunc    fs.WalkDirFunc
	FileCounter int64
	wg          sync.WaitGroup
	isCompare   bool
}

type Producer func() scanner.Action

func (w *Walker) ParWalk(workDir string) {
	w.wg.Add(1)
	go func() {
		if err := filepath.WalkDir(workDir, w.walkFunc); err != nil {
			fmt.Println(err)
		}
		w.wg.Done()
	}()

	var wg2 sync.WaitGroup
	if w.isCompare {
		wg2.Add(1)
		StartCompare(&wg2)
	}

	w.wg.Wait()
	close(parsed)
	wg2.Wait()
}

func NewWalker(actionMap map[string][]Producer, mode scanner.OpMode, startCompare bool) (w *Walker) {
	w = &Walker{}
	var resultHandler func(p base.Case) = nil
	if startCompare {
		resultHandler = FeedCase
	}
	parsed = make(chan base.Case, 8)
	w.walkFunc = func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), "~") && (strings.HasSuffix(d.Name(), ".xlsm") || strings.HasSuffix(d.Name(), ".xlsx")) {
			w.wg.Add(1)
			go func(filePath string) {
				asc := scanner.New(filePath, mode)
				if asc.Valid() {
					atomic.AddInt64(&w.FileCounter, 1)
					for s, producers := range actionMap {
						for _, producer := range producers {
							asc.RegisterAction(s, producer())
						}
					}
					asc.Scan(resultHandler)
				} else {
					log.Printf("file:%s process failed\n", filePath)
				}
				w.wg.Done()
			}(path)
		}
		return nil
	}
	return
}
