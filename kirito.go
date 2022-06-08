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

func Init(run func(), shutdownSignal chan bool) {
	mu.Lock()
	running = true
	mu.Unlock()

	hasRun := make(chan bool)

	go func() {
		run()

		if shutdownSignal != nil {
			<-shutdownSignal
		}

		hasRun <- true
	}()

	for {
		select {
		case fun := <-functionQueue:
			fun()
		case <-hasRun:
			return
		default:
		}
	}
}

func Queue(fun func()) {
	functionQueue <- fun
}

func QueueBlocking(fun func()) {
	Get(func() any {
		fun()
		return nil // not like return type matters
	})
}

func Get[R any](fun func() R) R {
	return <-GetAsync(fun)
}

func GetAsync[R any](fun func() R) chan R {
	hasRun := make(chan R)

	functionQueue <- func() {
		hasRun <- fun()
	}

	return hasRun
}

func Running() bool {
	mu.Lock()
	defer mu.Unlock()

	return running
}
