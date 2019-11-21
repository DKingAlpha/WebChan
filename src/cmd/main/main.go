package main

import (
	"fmt"
	"internal/shared_vars"
	"internal/utils"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var queues *sync.Map = nil

func NewTimeoutQueueWithParam(channelId string, key string, perm utils.ChanPerm) *utils.TimeoutQueue {
	return utils.NewTimeoutQueue(channelId, 1000, 60*60*24*7, key, perm) // 7 days
}


func rootHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")[1:]
	q := utils.GetUrlArgs(r.URL.RawQuery)
	key, foundKey := q["key"]
	if !foundKey {
		key = ""
	}
	adminMode := false
	adminPassword, foundAdminPassword := q["admin"]
	if foundAdminPassword {
		adminMode = adminPassword == shared_vars.AdminPassword
	}

	switch r.Method {
	case "POST":
		// write data to channel
		if len(p) < 2 || p[1] == "" {
			_, _ = fmt.Fprintln(w, "request url format: POST /channel/msg")
			return
		}
		permS, foundPerm := q["perm"]
		if foundPerm {
			permS = ""
		}
		perm := *utils.GetPerm(permS)

		queue, foundPerm := queues.Load(p[0])
		if !foundPerm {
			queue = NewTimeoutQueueWithParam(p[0], key, perm)
			queues.Store(p[0], queue)
		}
		// check owner
		if key == queue.(*utils.TimeoutQueue).Key {
			if foundPerm {
				// owner is updating perm
				queue.(*utils.TimeoutQueue).Perm = perm
			}
		} else {
			if !adminMode && queue.(*utils.TimeoutQueue).Perm.W == 0 {
				_, _ = fmt.Fprintln(w, "Wrong key to channel")
				return
			}
		}
		chanMsg := &utils.ChanMessage{
			Ch: p[0],
			T: shared_vars.CurrentTime,
			M: p[1],
		}
		_, _ = fmt.Fprintln(w, queue.(*utils.TimeoutQueue).Enqueue(chanMsg))
	case "GET":
		// show data in channel
		if len(p) < 1 {
			_, _ = fmt.Fprintln(w, "request url format: GET /channel | GET /channel/from | GET /channel/from/to")
			return
		}
		if len(p) >= 1 {
			var from int64 = 0
			var to int64 = math.MaxInt64
			var err error = nil
			if len(p) >= 2 && p[1] != "" {
				from, err = strconv.ParseInt(p[1], 10, 64)
				if err != nil {
					_, _ = fmt.Fprintln(w, "wrong timestamp: from")
					return
				}
			}
			if len(p) >= 3  && p[2] != "" {
				to, err = strconv.ParseInt(p[2], 10, 64)
				if err != nil {
					_, _ = fmt.Fprintln(w, "wrong timestamp: to")
					return
				}
			}
			queue, found := queues.Load(p[0])
			if found {
				if !adminMode && key != queue.(*utils.TimeoutQueue).Key && queue.(*utils.TimeoutQueue).Perm.R == 0 {
					_, _ = fmt.Fprintln(w, "Wrong key to channel")
					return
				}
				switch len(p) {
				case 2:
					_, _ = fmt.Fprint(w, queue.(*utils.TimeoutQueue).GetDataFrom(from))
				case 3:
					_, _ = fmt.Fprint(w, queue.(*utils.TimeoutQueue).GetDataFromTo(from, to))
				default:
					_, _ = fmt.Fprint(w, queue.(*utils.TimeoutQueue).GetData())
				}
			} else {
				_, _ = fmt.Fprint(w, "")
			}
		}
	case "DELETE":
		if len(p) < 1 || p[0] == "" {
			_, _ = fmt.Fprintln(w, "request url format: DELETE /channel")
			break
		}
		if !adminMode {
			if queue, found := queues.Load(p[0]); found && key != queue.(*utils.TimeoutQueue).Key &&
				queue.(*utils.TimeoutQueue).Perm.D == 0{
				_, _ = fmt.Fprintln(w, "Wrong key to channel")
				return
			}
		}
		queues.Delete(p[0])
		_, _ = fmt.Fprintln(w, "OK")
	default:
		_, _ = fmt.Fprintln(w, "Unsupported HTTP Method: " + r.Method)
	}
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

	queues = utils.LoadSyncMap(shared_vars.DumpJSONPath)
	if queues == nil {
		queues = &sync.Map{}
	}

	go signalCatcher()
	go queuesCleaner()
	go statusReporter()

	addr := os.Args[1]
	http.HandleFunc("/", limitClient(rootHandler, 200))
	log.Fatal(http.ListenAndServe(addr, nil))
}

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
		case <-time.After(time.Minute):
			log.Printf("Active Channel: %d\n", syncMapCount(queues))
		}
	}
}


func queuesCleaner() {
	for {
		queues.Range(func(key, queue interface{}) bool {
			select {
			// clean one channel every seconds
			case <- time.After(time.Second / 10):
				shared_vars.CurrentTime = time.Now().Unix()
				queue.(*utils.TimeoutQueue).CleanTimeout()
				if queue.(*utils.TimeoutQueue).Empty() {
					queues.Delete(key)
				}
			}
			return true
		})
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
	utils.DumpSyncMap(queues, shared_vars.DumpJSONPath)
	os.Exit(0)
}
