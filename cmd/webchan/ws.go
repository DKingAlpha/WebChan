package main

import (
	"github.com/DKingCN/WebChan/internal/utils"
	"golang.org/x/net/websocket"
	"strings"
	"sync"
	"time"
)

var QueueWSNotifyChan = make(chan *utils.ChanMessage, 128)
var queueWSNotify = sync.Map{}		// synced type: map[string]map[*websocket.Conn]bool


func wsHandler(ws *websocket.Conn) {
	var channel string
	if paths := strings.Split(ws.Config().Location.Path, "/"); len(paths) >= 3 {
		channel = paths[2]
	} else {
		_ = ws.Close()
		return
	}
	subscribers, _ := queueWSNotify.LoadOrStore(channel, &sync.Map{})
	subscribers.(*sync.Map).Store(ws, true)

	for {
		time.Sleep(5 * time.Second)
		subs, ok := queueWSNotify.Load(channel)
		if !ok {
			return
		}
		_, ok = subs.(*sync.Map).Load(ws)
		if !ok {
			return
		}
	}

}

func limitWSClient(handler websocket.Handler, clientNum int) websocket.Handler {
	sema := make(chan bool, clientNum)
	return func(ws *websocket.Conn) {
		sema <- true
		defer func() {<- sema }()
		handler(ws)
	}
}
