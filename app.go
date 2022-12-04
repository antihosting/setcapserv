/*
 *
 * Copyright 2022-present Zander Schwid & Co. LLC. All rights reserved.
 *
 */

package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	done chan struct{}
	err atomic.Error
}

func NewApp() *App {
	t := &App{
		done: make(chan struct{}),
	}
	t.init()
	return t
}

func (t *App) Context() context.Context {
	return t
}

func (t *App) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (t *App) Done() <-chan struct{} {
	return t.done
}

func (t *App) Err() error {
	return t.err.Load()
}

func (t *App) Value(key interface{}) interface{} {
	return nil
}

func (t *App) init() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for signal := range signalChan {
			fmt.Printf("Signal '%v'\n", signal)
			t.err.Store(errors.Errorf("interrupted by signal '%v'", signal))
			t.done <- struct{}{}
			break
		}
	}()
}
