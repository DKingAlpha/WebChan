package main

import (
	"fmt"
	"github.com/DKingCN/WebChan/internal/shared_vars"
	"golang.org/x/net/websocket"
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
		time.Sleep(5*time.Minute)
		log.Printf("Active Channel: %d\n", syncMapCount(queues))
	}
}


func queuesCleaner() {
	for {
		time.Sleep(time.Second)
		queues.Range(func(key, queue interface{}) bool {
			time.Sleep(time.Second)
			shared_vars.CurrentTime = time.Now().Unix()
			queue.(*RTQ).CleanTimeout()
			if queue.(*RTQ).Empty() {
				queues.Delete(key)
			}
			return true
		})
	}
}

func activityCleaner() {
	for {
		time.Sleep(time.Second)
		shared_vars.CurrentTime = time.Now().Unix()
		activityLog.Clean()
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

func WSNotifyDaemon() {
	for {
		chanmsg := <- QueueWSNotifyChan
		subscribers, ok := queueWSNotify.Load(chanmsg.Ch)
		if !ok {
			continue
		}
		subscribers.(*sync.Map).Range(func(w interface{}, _ interface{}) bool {
			go func() {
				if err := websocket.Message.Send(w.(*websocket.Conn), fmt.Sprintf("%d|%s\n", chanmsg.T, chanmsg.M)); err != nil{
					_ = w.(*websocket.Conn).Close()
					subscribers.(*sync.Map).Delete(w)
				}
			}()
			return true
		})
	}
}
