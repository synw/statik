package watcher

import (
	"fmt"
	"sync"
	"time"

	wa "github.com/radovskyb/watcher"

	"github.com/synw/statik/conf"
	"github.com/synw/statik/ws"
)

var w = wa.New()
var building = false
var changeRequestedWhileBuilding = false

func Watch(verbose bool, reload bool) {
	for dir := range conf.WatchBuilders {
		if verbose {
			fmt.Println("Watching for changes in", dir)
		}
		err := w.AddRecursive(dir)
		if err != nil {
			panic("Can not add path lib")
		}
	}
	w.FilterOps(wa.Write, wa.Create, wa.Move, wa.Remove, wa.Rename)
	// lauch listener
	var mux sync.Mutex
	go func() {
		for {
			select {
			case e := <-w.Event:
				if building == false {
					go build(&mux, e.Path, verbose, reload)
				} else {
					if verbose {
						fmt.Println("Change requested while building, delaying next build")
					}
					changeRequestedWhileBuilding = true
				}
			case err := <-w.Error:
				msg := "Watcher error " + err.Error()
				fmt.Println(msg)
			case <-w.Closed:
				msg := "Watcher closed"
				fmt.Println(msg)
				return
			}
		}
	}()
	// start listening
	err := w.Start(time.Millisecond * 200)
	if err != nil {
		panic("Error starting the watcher")
	}
}

func build(mux *sync.Mutex, path string, verbose bool, reload bool) {
	building = true
	msg := "Change detected in " + path
	if verbose {
		fmt.Println(msg)
	}
	if verbose {
		fmt.Println("Running build ...")
	}
	mux.Lock()
	runBuildCmd(path, verbose)
	mux.Unlock()
	if verbose {
		fmt.Println("Build done, reloading")
	}
	if reload {
		ws.SendMsg(msg)
	}
	building = false
	if changeRequestedWhileBuilding == true {
		changeRequestedWhileBuilding = false
		build(mux, path, verbose, reload)
	}
}
