package main

import (
	"encoding/json"
	"internal/shared_vars"
	"internal/utils"
	"io/ioutil"
	"log"
	"os"
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

func (p* ChanPerm) String() string {
	perm := ""
	if p.R != 0 {
		perm += "r"
	} else {
		perm += "-"
	}
	if p.W != 0 {
		perm += "w"
	} else {
		perm += "-"
	}
	if p.D != 0 {
		perm += "d"
	} else {
		perm += "-"
	}
	return perm
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

func (al * ActivityLog) Remove(channelId string) {
	al.lock.Lock()
	defer al.lock.Unlock()
	delete(al.Acts, channelId)
}


func (al *ActivityLog) Log(channelId string) bool {
	al.lock.Lock()
	defer al.lock.Unlock()
	if len(al.Acts) >= al.Cap {
		// just return false, wait for the next clean-up
		return false
	}
	if act, found := al.Acts[channelId]; found {
		act.Count++
		act.LastTime = shared_vars.CurrentTime
	} else {
		al.Acts[channelId] = &Activity{
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

func (al *ActivityLog) Clean() {
	al.lock.Lock()
	defer al.lock.Unlock()
	if len(al.Acts) < al.Cap {
		return
	}

	flat := make([]ChanAct, len(al.Acts))
	i := 0
	for channel, act := range al.Acts {
		flat[i] = ChanAct{channel, act}
		i++
	}
	sort.Slice(flat, func(i, j int) bool {
		return flat[i].Act.LastTime > flat[j].Act.LastTime
	})
	// remove 10% oldest
	trimmed := flat[:len(flat)*9/10]
	al.Acts = map[string]*Activity{}
	for _, kv := range trimmed {
		al.Acts[kv.Chan] = kv.Act
	}
}

func (al *ActivityLog) Rank() *[]ChanAct{
	al.lock.Lock()
	defer al.lock.Unlock()
	flat := make([]ChanAct, len(al.Acts))
	i := 0
	for channel, act := range al.Acts {
		flat[i] = ChanAct{channel, act}
		i++
	}
	sort.Slice(flat, func(i, j int) bool {
		return flat[i].Act.Count >  flat[j].Act.Count
	})
	return &flat
}


func LoadActivityLog(path string) *ActivityLog {
	log.Println("Loading activity from ", path)
	f := ActivityLog{
		Cap:     100,
		Timeout: shared_vars.ActivityTimeout,
		Acts:    map[string]*Activity{},
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load activity from %s: %v\n", path, err)
		return &f
	}
	if err := json.Unmarshal(data, &f); err != nil {
		log.Printf("Failed to unmarshal activity: %v\n", err)
		return &f
	}
	log.Printf("Loaded %s\n", path)
	return &f
}

func (al *ActivityLog) Dump(path string) {
	log.Println("Dumping activity to ", path)
	data, err := json.Marshal(al)
	if err != nil {
		log.Printf("Failed to marshal activity: %v\n", err)
		return
	}
	if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
		log.Printf("Failed to dump activity to %s: %v\n", path, err)
	}
}
