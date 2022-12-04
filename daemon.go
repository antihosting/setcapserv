/*
 *
 * Copyright 2022-present Zander Schwid & Co. LLC. All rights reserved.
 *
 */

package main

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	targetPattern = "%1"
)

type watchDaemon struct {

	watchFilePath  string
	commands       [][]string

	log  *log.Logger
	verbose  bool

	last    atomic.Value // struct watchFileInfo

	watcher      *fsnotify.Watcher
	triggerTime  atomic.Int64
}

type watchFileInfo struct {
	Size  int64
	ModTime  time.Time
}

func NewDeamon(watchFilePath string, commands [][]string, log  *log.Logger, verbose  bool) *watchDaemon {
	return &watchDaemon {
		watchFilePath: watchFilePath,
		commands: commands,
		log: log,
		verbose: verbose,
	}
}

func (t *watchDaemon) Run(ctx context.Context) (err error) {

	if t.watcher != nil {
		return errors.New("watcher already running")
	}

	t.watchFilePath, err = filepath.Abs(t.watchFilePath)
	if err != nil {
		return err
	}

	if stat, err := os.Stat(t.watchFilePath); err == nil {
		t.last.Store(watchFileInfo{
			Size: stat.Size(),
			ModTime: stat.ModTime(),
		})
		if t.verbose {
			t.log.Printf("Watcher: File exist '%s' size='%d', modTime='%v'\n", t.watchFilePath, stat.Size(), stat.ModTime())
		}
	} else {
		// file not exist yet
		t.last.Store(watchFileInfo{
			Size: 0,
			ModTime: time.Unix(0, 0),
		})
		if t.verbose {
			t.log.Printf("Watcher: File not exist '%s'\n", t.watchFilePath)
		}
	}

	t.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}

	defer t.watcher.Close()

	err = t.watcher.Add(t.watchFilePath)
	if err != nil {
		t.log.Printf("Watcher: can not add file '%s' to watcher, %v", t.watchFilePath, err)
		return
	}

	watchFileDir := filepath.Dir(t.watchFilePath)
	err = t.watcher.Add(watchFileDir)
	if err != nil {
		t.log.Printf("Watcher: can not add path '%s' to watcher, %v", watchFileDir, err)
		return
	}

	for {
		select {
		case event, ok := <-t.watcher.Events:
			if ok {
				t.onEvent(event)
			}
		case err, ok := <-t.watcher.Errors:
			if ok {
				t.log.Printf("WatcherError: %v\n", err)
			}
		case <-ctx.Done():
			t.log.Printf("WatcherStopped: %v\n", ctx.Err())
			return nil
		}
	}


}

func (t *watchDaemon) onEvent(event fsnotify.Event) {

	if t.verbose {
		t.log.Printf("Watcher: event '%s'\n", event.String())
	}

	if event.Op == fsnotify.Create || event.Op == fsnotify.Write {

		if t.watchFilePath != event.Name {
			t.log.Printf("Watcher: nnexpected event '%s'\n", event.Name)
			return
		}
		stat, err := os.Stat(event.Name)
		if err != nil {
			t.log.Printf("Watcher: file not found '%s', %v\n", event.Name, err)
			return
		}

		ls, ok := t.last.Load().(watchFileInfo)
		if !ok {
			t.log.Println("Watcher: invalid value in last stat")
			return
		}

		if ls.Size != stat.Size() || ls.ModTime != stat.ModTime() {

			t.last.Store(watchFileInfo{
				Size: stat.Size(),
				ModTime: stat.ModTime(),
			})

			t.log.Printf("Watcher: file changed '%s'\n", event.Name)
			t.trigger()

		}

	}

}

func (t *watchDaemon) trigger() {
	current := time.Now().UnixNano()
	t.triggerTime.Store(current)
	time.AfterFunc(10 * time.Millisecond, func() {
		if t.triggerTime.Load() == current {
			if isFileLocked(t.watchFilePath) {
				// try to update again
				t.log.Printf("Watcher: file locked '%s'\n", t.watchFilePath)
				t.trigger()
			} else {
				t.deploy()
			}
		}
	})
}

func (t *watchDaemon) deploy() {

	for _, command := range t.commands {

		if len(command) < 1 {
			t.log.Printf("Invalid command '%v'\n", command)
			continue
		}

		var args []string
		for _, value := range command {
			if value == targetPattern {
				args = append(args, t.watchFilePath)
			} else {
				args = append(args, value)
			}
		}

		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.log.Printf("RunCommand: error '%v', output '%s', %v\n", cmd, output, err)
			return
		}

		t.log.Printf("RunCommand: success '%v', output '%s'\n", cmd, output)

	}

}

func isFileLocked(filePath string) bool {
	if file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_EXCL, 0); err != nil {
		return true
	} else {
		file.Close()
		return false
	}
}
