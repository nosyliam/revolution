package config

import (
	"context"
	"fmt"
	"github.com/sqweek/dialog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
)

type event struct {
	path, op string
	value    interface{}
}

type Runtime struct {
	ready       bool
	roots       map[string]Reactive
	events      []event
	errorActive bool
}

func (r *Runtime) handleError(op string, err error) bool {
	if err == nil {
		return true
	}
	if r.errorActive {
		return false
	}
	r.errorActive = true
	dialog.Message(fmt.Sprintf("%s operation failed: %v", op, err)).Error()
	r.errorActive = false
	return false
}

func (r *Runtime) AddRoot(name string, object Reactive) {
	r.roots[name] = object
}

func (r *Runtime) Set(path string, value interface{}) {
	if !r.ready {
		r.events = append(r.events, event{path: path, op: "set", value: value})
		return
	}
	runtime.EventsEmit(AppContext, "set", path, value)
}

func (r *Runtime) Append(path string, primitive bool) {
	if !r.ready {
		r.events = append(r.events, event{path: path, op: "append"})
		return
	}
	runtime.EventsEmit(AppContext, "append", path, primitive)
}

func (r *Runtime) Delete(path string) {
	if !r.ready {
		r.events = append(r.events, event{path: path, op: "delete"})
		return
	}
	runtime.EventsEmit(AppContext, "delete", path)
}

func (r *Runtime) Listen() {
	_ = runtime.EventsOn(AppContext, "set_client", func(data ...interface{}) {
		path := data[0].(string)
		if r.handleError(
			"Set",
			r.roots[getRoot(path)].Set(mustCompilePath(getPath(path)), 0, data[1]),
		) {
			r.handleError("Save", r.roots[getRoot(path)].File().Save())
		}
	})
	_ = runtime.EventsOn(AppContext, "append_client", func(data ...interface{}) {
		path := data[0].(string)
		if r.handleError(
			"Append",
			r.roots[getRoot(path)].Append(mustCompilePath(getPath(path)), 0),
		) {
			r.handleError("Save", r.roots[getRoot(path)].File().Save())
		}
	})
	_ = runtime.EventsOn(AppContext, "delete_client", func(data ...interface{}) {
		path := data[0].(string)
		if r.handleError(
			"Delete",
			r.roots[getRoot(path)].Append(mustCompilePath(getPath(path)), 0),
		) {
			r.handleError("Save", r.roots[getRoot(path)].File().Save())
		}
	})
}

func (r *Runtime) Start() {
	r.Listen()
	for _, evt := range r.events {
		switch evt.op {
		case "set":
			runtime.EventsEmit(AppContext, "set", evt.path, evt.value)
		case "append":
			runtime.EventsEmit(AppContext, "append", evt.path)
		}
	}
	r.events = nil
}

func NewRuntime(ctx context.Context) *Runtime {
	AppContext = ctx
	app := &Runtime{roots: make(map[string]Reactive)}
	ready := make(chan struct{})
	runtime.EventsOnce(ctx, "ready", func(...interface{}) {
		app.ready = true
		for len(app.roots) != 3 {
			fmt.Println(len(app.roots))
			<-time.After(100 * time.Millisecond)
		}
		ready <- struct{}{}
	})
	go func() {
		select {
		case <-ready:
			app.Start()
			return
		case <-time.After(time.Second * 10):
			dialog.Message("UI failed to start! Please contact the developer for assistance").Error()
			//os.Exit(1)
		}
	}()
	return app
}
