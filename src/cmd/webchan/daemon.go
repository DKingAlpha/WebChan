package main

import (
	"internal/shared_vars"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func statusReporter() {
	syncMapCount := func(p *sync.Map) int {
		count := 0
		p.Range(func(key, _ interface{}) bool {
			count += 1
			return true
		})
		return count
	}
	for {
		select {
		case <-time.After(5*time.Minute):
			log.Printf("Active Channel: %d\n", syncMapCount(queues))
		}
	}
}


func queuesCleaner() {
	for {
		queues.Range(func(key, queue interface{}) bool {
			select {
			// clean one channel every 100ms
			case <- time.After(time.Second / 10):
				shared_vars.CurrentTime = time.Now().Unix()
				queue.(*RTQ).CleanTimeout()
				if queue.(*RTQ).Empty() {
					queues.Delete(key)
				}
			}
			return true
		})
	}
}

func activityCleaner() {
	for {
		select {
		// clean one channel every 10 seconds
		case <- time.After(time.Second*10):
			shared_vars.CurrentTime = time.Now().Unix()
			activityLog.Clean()
		}
	}
}


func signalCatcher() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)
	s := <-sigc
	log.Println("Caught signal: ", s.String())
	DumpSyncMap(queues, shared_vars.DumpDBPath)
	activityLog.Dump(shared_vars.DumpActivityPath)
	os.Exit(0)
}
