package walker

import (
	"bufio"
	"bytes"
	"errors"
	"extool/action"
	"extool/action/approximate"
	"extool/config"
	"extool/module"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type OpMode int

const (
	READONLY OpMode = iota
	WRITESOURCE
	WRITECOPY
)

type Walker struct {
	actionMap    map[string][]func() action.Action
	listenerMap  map[string][]func() action.Listener
	caseComparer *approximate.CaseComparer
	dirFunc      fs.WalkDirFunc
	mode         OpMode
	workDir      string
	callbacks    []func()
	saveAs       string
	startTime    time.Time
	fileCount    int64
	dirCount     int64
	walkedCount  int64
	failedCount  int64
}

func NewWalker(mode OpMode) *Walker {
	var walker = &Walker{
		listenerMap: make(map[string][]func() action.Listener, 4),
		actionMap:   make(map[string][]func() action.Action, 4),
		mode:        mode,
	}
	walker.startTime = time.Now()
	walker.dirFunc = func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !strings.HasPrefix(d.Name(), "~") && (strings.HasSuffix(d.Name(), ".xlsm") || strings.HasSuffix(d.Name(), ".xlsx")) {
			atomic.AddInt64(&walker.walkedCount, 1)
			// if everything ok,open the file
			file, err := excelize.OpenFile(path)
			if err != nil {
				atomic.AddInt64(&walker.failedCount, 1)
				log.Println(err)
				return nil
			}
			// pose action per sheet
			for sheetName, producers := range walker.actionMap {
				if len(producers) == 0 {
					continue
				}
				// prepare context for every action pose on same sheet
				ctx := action.Context{File: file, Scase: module.NewCase(path), IsReaOnly: walker.mode == READONLY}
				ctx.RowNum = 0
				ctx.SheetName = sheetName // populate sheet name
				rows, err := file.Rows(sheetName)
				if err != nil {
					log.Println(err)
				}
				var actions = make([]action.Action, 0, 4)
				for _, producer := range producers {
					actions = append(actions, producer())
				}
				for rows.Next() {
					ctx.RowNum++ // populate current row number
					row, err := rows.Columns()
					if err != nil {
						fmt.Println(err)
					}
					ctx.Row = row // populate current row info
					// just do it
					for _, ac := range actions {
						ac(&ctx)
					}
				}
				for _, listener := range walker.listenerMap[sheetName] {
					listener()(&ctx)
				}
			}
			// after all action done,do different kind of operation depends on the mode
			switch walker.mode {
			case READONLY:
			case WRITESOURCE:
				err := file.Save()
				if err != nil {
					log.Println(err)
				}
			case WRITECOPY:
				rel, _ := filepath.Rel(walker.workDir, path)
				errs := file.SaveAs(filepath.Join(walker.saveAs, rel))
				if errs != nil {
					if errors.Is(errs, syscall.ENOTDIR) {
						rel, err2 := filepath.Rel(walker.workDir, filepath.Dir(path))
						if err2 != nil {
							log.Println(err2)
						}
						err := os.MkdirAll(filepath.Join(walker.saveAs, rel), os.ModeDir)
						if err != nil && !errors.Is(err, syscall.ERROR_ALREADY_EXISTS) {
							log.Println(err)
						}
					}
				}
			}
		} else if d.IsDir() {
			atomic.AddInt64(&walker.dirCount, 1)
			return nil
		}
		atomic.AddInt64(&walker.fileCount, 1)
		return nil
	}
	return walker
}

func (w *Walker) RegisterAction(sheet string, action func() action.Action) *Walker {
	w.actionMap[sheet] = append(w.actionMap[sheet], action)
	return w
}

func (w *Walker) RegisterListener(sheet string, listener func() action.Listener) *Walker {
	w.listenerMap[sheet] = append(w.listenerMap[sheet], listener)
	return w
}

// register a callback function,only be invoked when walker finish his walk

func (w *Walker) RegisterCallBack(callback func()) *Walker {
	w.callbacks = append(w.callbacks, callback)
	return w
}

func (w *Walker) WithCaseCompare() *Walker {
	w.caseComparer = approximate.NewCaseComparer()
	CaseFeedListener := w.caseComparer.GetCaseFeedListener()
	w.RegisterAction(module.Sheet1, approximate.CaseParseAction).
		RegisterListener(module.Sheet1, CaseFeedListener).
		RegisterCallBack(w.caseComparer.Close)

	return w
}

func (w *Walker) inheritSource() {
	var inheritSource = config.Cfg.InheritSource

	if len(inheritSource) > 0 {
		if len(inheritSource) > 8 {
			log.Fatalln("Only support inheritance report within 8 files")
		}
		for i, s := range inheritSource {
			f, err := os.Open(s)
			if err != nil {
				log.Println(err)
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				ucase := module.UnMarshal([]byte(scanner.Text()))
				if i == 0 {
					w.caseComparer.Chain.PushBack(ucase)
				} else {
					go func() {
						w.caseComparer.Parsed <- ucase
					}()
				}
			}
		}
	}
}

// only return after all walkFunc finish

func (w *Walker) GoWalkDir(workDir string) *Walker {
	w.workDir = workDir
	np := filepath.Join(filepath.Dir(w.workDir), "copy")
	errm := os.Mkdir(np, os.ModeDir)
	if errm != nil {
		if errors.Is(errm, syscall.ERROR_ALREADY_EXISTS) {
			// Noop
		} else {
			log.Println(errm)
		}
	}
	w.saveAs = np

	dirFunc, waiter := wrappedWalkFunc(w.dirFunc)

	var finish = func() {}
	if w.caseComparer != nil {
		w.inheritSource()
		finish = w.caseComparer.StartCompareBackGround()
	}

	err := filepath.WalkDir(workDir, dirFunc)
	if err != nil {
		log.Println(err)
	}

	waiter()
	for _, callback := range w.callbacks {
		callback()
	}
	finish()
	return w
}

func (w *Walker) Report() {
	log.Printf("All task completed\n"+
		"  <Dir Count:%d>\n"+
		"  <Total File Count:%d>\n"+
		"  <Handled File Count:%d>\n"+
		"  <Failed File Count:%d>\n"+
		"  <Remaining File Count:%d>\n"+
		"  <Spend Time:%s>\n"+
		"Taking a while to generate the report\n",
		w.dirCount, w.fileCount, w.walkedCount, w.failedCount, w.caseComparer.Chain.Len(), time.Since(w.startTime))
	if w.caseComparer != nil {
		var location string
		switch w.mode {
		case WRITECOPY:
			location = filepath.Join(w.saveAs, "report")
		default:
			location = filepath.Join(w.workDir, "report")
		}
		var (
			filename         = "[" + time.Now().Format("2006-01-02 ※ 15-04-05") + "]" + ".txt"
			absoluteFilePath = filepath.Join(location, filename)
		)
		errm := os.Mkdir(location, os.ModeDir)
		if errm != nil && !errors.Is(errm, fs.ErrExist) {
			log.Printf("Failed to create directory %s", location)
		}
		f, erro := os.OpenFile(absoluteFilePath, os.O_RDWR|os.O_CREATE, 0666)
		if erro != nil {
			log.Printf("Failed to create report:%s\n", filename)
		}
		defer func() {
			err := f.Close()
			log.Println(err)
		}()
		bufc := make(chan *bytes.Buffer, 16)
		wg := sync.WaitGroup{}
		callback := func() {
			wg.Done()
		}
		for ele := w.caseComparer.Chain.Front(); ele != nil; ele = ele.Next() {
			wg.Add(1)
			go ele.Value.(module.Case).Marshal(bufc, callback)
		}
		go func() {
			wg.Wait()
			close(bufc)
		}()

		for buf := range bufc {
			_, err := f.Write([]byte(buf.String()))
			if err != nil {
				log.Println(err)
			}
		}

		log.Printf("Detailed report located at %s\n", absoluteFilePath)
	}
}

// wrap the dirFunc,provide following mechanism
// - parallel execute the providing dirFunc
// - limit the dirFunc call within 8 the same time,otherwise stuck the DirWalker
// - provide a waiter for all wrapped dirFunc call,waiter will return after all wrapped dirFunc finish

func wrappedWalkFunc(dirFunc fs.WalkDirFunc) (fs.WalkDirFunc, func()) {
	var wg sync.WaitGroup
	var token = make(chan struct{}, 8)
	return func(path string, d fs.DirEntry, err error) error {
			wg.Add(1)
			token <- struct{}{}
			go func() {
				err := dirFunc(path, d, err)
				if err != nil {
					log.Println(err)
				}
				<-token
				wg.Done()
			}()
			return nil
		},
		func() {
			wg.Wait()
		}
}
