package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"internal/shared_vars"
	"internal/utils"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")[1:]
	q := GetUrlArgs(r.URL.RawQuery)

	// get runtime method
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
		bodyb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			bodyb = []byte("")
		}
		body := string(bodyb)
		postUsage := "request url format:\nPOST /channel/msg\n\nPOST /channel\nDATA\n\nPOST /channel/data1\nDATA2"

		msg := ""
		// try append msg in url
		if len(p) >= 2 && p[1] != "" {
			msg += p[1]
		}
		if body != "" {
			msg += "\n" + body
		}

		validRequest := false
		permS, foundPerm := q["perm"]

		// create new channel if not exists
		queue, foundQueue := queues.Load(p[0])
		if !foundQueue {
			// create new channel with proper permission
			var perm *ChanPerm = nil
			if foundKey {
				// private, perm default 0,0,0
				if foundPerm {
					perm = GetPerm(permS)
				} else {
					perm = &ChanPerm{}
				}
			} else {
				// public, perm default 1,1,1
				perm = &ChanPerm{1,1,1}
			}
			queue = NewRTQWithParam(p[0], key, perm)
			queues.Store(p[0], queue)
			validRequest = true
		} else {
			// channel exists

			// everyone could change property of public channel. consider protecting it by key=XXX&perm=r
			if queue.(*RTQ).Key == "" || queue.(*RTQ).Key == key {
				// public channel or owner of private channel
				if foundPerm {
					queue.(*RTQ).Perm = *GetPerm(permS)
					validRequest = true
				}
				if foundKey {
					queue.(*RTQ).Key = key
				}
			} else {
				if !adminMode && queue.(*RTQ).Perm.W == 0 {
					_, _ = fmt.Fprintln(w, "Private channel")
					return
				}
			}
		}
		// still not matching a valid update request. the only operation left is to push msg and its nil
		if !validRequest && msg == "" {
			_, _ = fmt.Fprintln(w, postUsage)
			return
		}

		// log this to public log only when someone else could read
		if queue.(*RTQ).Perm.R != 0 {
			activityLog.Log(p[0])
		}
		chanMsg := &utils.ChanMessage{
			Ch: p[0],
			T: shared_vars.CurrentTime,
			M: msg,
		}
		_, _ = fmt.Fprintln(w, queue.(*RTQ).Enqueue(chanMsg))
	case "GET":
		// show data in channel
		if len(p) == 0 || p[0] == "" {
			statusHandler(w, r)
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
				if !adminMode {
					if queue.(*RTQ).Perm.R == 0 {
						if queue.(*RTQ).Key == "" {
							_, _ = fmt.Fprintln(w, "No read permission on this public channel")
							return
						}
						if queue.(*RTQ).Key != key {
							_, _ = fmt.Fprintln(w, "Private channel")
							return
						}
					}
				}
				_, showTime := q["time"]
				switch len(p) {
				case 2:
					_, _ = fmt.Fprint(w, queue.(*RTQ).GetDataFrom(from, showTime))
				case 3:
					_, _ = fmt.Fprint(w, queue.(*RTQ).GetDataFromTo(from, to, showTime))
				default:
					_, _ = fmt.Fprint(w, queue.(*RTQ).GetData(showTime))
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
			if queue, found := queues.Load(p[0]); found && queue.(*RTQ).Perm.D == 0 {
				if queue.(*RTQ).Key == "" {
					_, _ = fmt.Fprintln(w, "No delete permission on this public channel")
					return
				}
				if queue.(*RTQ).Key != key {
					_, _ = fmt.Fprintln(w, "Private channel")
					return
				}
			}
		}
		queues.Delete(p[0])
		activityLog.Remove(p[0])
		_, _ = fmt.Fprintln(w, "OK")
	default:
		_, _ = fmt.Fprintln(w, "Unsupported HTTP Method: " + r.Method)
	}
}
