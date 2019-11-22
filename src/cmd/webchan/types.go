package main

import (
	"internal/shared_vars"
	"internal/utils"
	"sort"
	"strings"
	"sync"
	"time"
)

type ChanPerm struct {
	// using int instead of bool to save space in JSON
	R int
	W int
	D int
}


func GetPerm(perm string) *ChanPerm {
	cp := ChanPerm{
		R: 0,
		W: 0,
		D: 0,
	}
	for _, p := range perm {
		if p == 'r' || p == 'R' {
			cp.R = 1
		}
		if p == 'w' || p == 'W' {
			cp.W = 1
		}
		if p == 'd' || p == 'D' {
			cp.D = 1
		}
	}
	return &cp
}

func GetUrlArgs(args string) map[string]string {
	q := map[string]string{}
	for _, kv :=  range strings.Split(args, "&") {
		kva := strings.SplitN(kv, "=", 2)
		v := ""
		if len(kva) == 2 {
			v = kva[1]
		}
		q[kva[0]] = v
	}
	return q
}


// RestrictedTimeoutQueue
type RTQ struct {
	utils.TimeoutQueue
	Key     string
	Perm    ChanPerm
}

func NewRTQ(channelId string, cap int, timeout int64, key string, perm ChanPerm) *RTQ {
	return &RTQ{
		TimeoutQueue: utils.TimeoutQueue{
			ChanId:  channelId,
			Cap:     cap,
			Timeout: timeout,
			Msgs:    nil,
		},
		Key:          key,
		Perm:         perm,
	}
}


type Activity struct {
	Count    int64
	LastTime int64
}

func (a Activity) GetTime() string {
	deltaT := time.Unix(shared_vars.CurrentTime, 0).Sub(time.Unix(a.LastTime, 0))
	return deltaT.String()
}

type ActivityLog struct {
	lock    sync.Mutex
	Cap     int
	Timeout int64
	Acts    map[string]*Activity
}



func (tq *ActivityLog) Log(channelId string) bool {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if len(tq.Acts) >= tq.Cap {
		// just return false, wait for the next clean-up
		return false
	}
	if act, found := tq.Acts[channelId]; found {
		act.Count++
		act.LastTime = shared_vars.CurrentTime
	} else {
		tq.Acts[channelId] = &Activity{
			Count:    1,
			LastTime: shared_vars.CurrentTime,
		}
	}
	return true
}

type ChanAct struct {
	Chan string
	Act  *Activity
}

func (tq *ActivityLog) Clean() {
	tq.lock.Lock()
	defer tq.lock.Unlock()
	if len(tq.Acts) < tq.Cap {
		return
	}

	flat := make([]ChanAct, len(tq.Acts))
	i := 0
	for channel, act := range tq.Acts {
		flat[i] = ChanAct{channel, act}
	}
	sort.Slice(flat, func(i, j int) bool {
		return flat[i].Act.LastTime < flat[j].Act.LastTime
	})
	// remove 10% oldest
	trimmed := flat[:len(flat)*9/10]
	tq.Acts = map[string]*Activity{}
	for _, kv := range trimmed {
		tq.Acts[kv.Chan] = kv.Act
	}
}

func (tq *ActivityLog) Rank() *[]ChanAct{
	tq.lock.Lock()
	defer tq.lock.Unlock()
	flat := make([]ChanAct, len(tq.Acts))
	i := 0
	for channel, act := range tq.Acts {
		flat[i] = ChanAct{channel, act}
	}
	sort.Slice(flat, func(i, j int) bool {
		return flat[i].Act.Count <  flat[j].Act.Count
	})
	return &flat
}
