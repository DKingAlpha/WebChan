package main

import (
	"internal/shared_vars"
	"log"
	"net/http"
	"os"
	"sync"
)

var queues *sync.Map = nil


func GetTimeoutBasedOnChannelId(channelId string) int64 {
	factor := len(channelId)
	if factor > 180 {
		factor = 180
	}
	return int64(60*60*24*factor)
}


func NewRTQWithParam(channelId string, key string, perm *ChanPerm) *RTQ {
	return NewRTQ(channelId, 1000, GetTimeoutBasedOnChannelId(channelId), key, *perm)
}


func limitClient(handler http.HandlerFunc, clientNum int) http.HandlerFunc {
	sema := make(chan bool, clientNum)
	return func(w http.ResponseWriter, req *http.Request) {
		sema <- true
		defer func() {<- sema }()
		handler(w, req)
	}
}


func main() {
	if len(os.Args) != 2 {
		log.Fatalln("missing argument addr:port")
	}

	queues = LoadSyncMap(shared_vars.DumpDBPath)
	activityLog = LoadActivityLog(shared_vars.DumpActivityPath)

	go signalCatcher()
	go queuesCleaner()
	go activityCleaner()
	go statusReporter()

	addr := os.Args[1]
	http.HandleFunc("/", limitClient(rootHandler, 200))
	http.HandleFunc("/tool", limitClient(toolHandler, 200))
	http.HandleFunc("/tool/", limitClient(toolHandler, 200))
	log.Fatal(http.ListenAndServe(addr, nil))
}
