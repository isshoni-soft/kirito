package kirito

import (
	"runtime"
	"sync"
)

var mu sync.Mutex

var functionQueue = make(chan func(), 20)
var running = false

func init() {
	runtime.LockOSThread()
}

func Run(run func()) {
	mu.Lock()
	running = true
	mu.Unlock()

	hasRun := make(chan bool)

	go func() {
		run()
		hasRun <- true
	}()

	for {
		select {
		case fun := <- functionQueue:
			fun()
		case <- hasRun:
			return
		default:
		}
	}
}

func Queue(fun func()) {
	functionQueue <- fun
}

func QueueBlocking(fun func()) {
	Get(func() interface {} {
		fun()
		return nil // not like return type matters
	})
}

func Get(fun func() interface {}) interface{} {
	hasRun := make(chan interface {})

	functionQueue <- func() {
		hasRun <- fun()
	}

	return <- hasRun
}

func Running() bool {
	mu.Lock()
	defer mu.Unlock()

	return running
}