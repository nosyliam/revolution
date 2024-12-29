package pattern

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/ast"
	"github.com/yuin/gopher-lua/parse"
)

type Metadata struct {
	Name string

	Position, Length, Width, Distance int
	RotateDirection, RotateCount      int
	InvertFB, InvertLR                bool

	BackpackPercentage int
	Minutes            int
	ShiftLock          bool
	ReturnMethod       string
	DriftComp          bool
	GatherPattern      bool
}

func NewMetadataFromStatementList(stmts []ast.Stmt) Metadata {
	meta := Metadata{}
	return meta
}

type Pattern struct {
	Meta  Metadata
	Proto *lua.FunctionProto

	Name string
	Path string
}

type Loader struct {
	watcher  *fsnotify.Watcher
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.RWMutex
	patterns map[string]*Pattern
}

func NewLoader() *Loader {
	w, _ := fsnotify.NewWatcher()
	ctx, cancel := context.WithCancel(context.Background())
	return &Loader{watcher: w, ctx: ctx, cancel: cancel, patterns: make(map[string]*Pattern)}
}

func (l *Loader) Start() error {
	os.MkdirAll("patterns", 0755)
	err := l.watcher.Add("patterns")
	if err != nil {
		return err
	}
	go l.run()
	return nil
}

func (l *Loader) run() {
	for {
		select {
		case event, ok := <-l.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				l.handleFileChange(event.Name)
			}
		case <-l.ctx.Done():
			return
		}
	}
}

func (l *Loader) handleFileChange(path string) {
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		return
	}
	if strings.HasSuffix(strings.ToLower(filepath.Ext(path)), ".lua") {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return
		}
		chunk, err := parse.Parse(bytes.NewReader(data), path)
		if err != nil {
			return
		}
		proto, err := lua.Compile(chunk, path)
		if err != nil {
			return
		}
		l.mu.Lock()
		//l.patterns[patternName] = proto
		l.mu.Unlock()
	}
}

func (l *Loader) Retrieve(pattern string) (*Pattern, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	p, ok := l.patterns[pattern]
	if !ok {
		return nil, fmt.Errorf("pattern not found: %s", pattern)
	}
	return p, nil
}

func (l *Loader) Close() error {
	l.cancel()
	return l.watcher.Close()
}
