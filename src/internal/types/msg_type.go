package types

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type ChanMessage struct {
	Ch string
	T int64
	M string
}

type Message struct {
	T int64
	M string
}

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
		if len(kva) != 2 {
			continue
		}
		q[kva[0]] = kva[1]
	}
	return q
}


func LoadSyncMap(path string) *sync.Map {
	log.Println("Loading data from ", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load file from %s : %v\n", path, err)
		return nil
	}
	tmpMap := map[string]*TimeoutQueue{}
	if err := json.Unmarshal(data, &tmpMap); err != nil {
		log.Printf("Failed to unmarshal data from %s : %v\n", path, err)
		return nil
	}
	f := sync.Map{}
	for key, value := range tmpMap {
		f.Store(key, value)
	}
	log.Printf("Loaded %s\n", path)
	return &f
}

func DumpSyncMap(p *sync.Map, path string) {
	log.Println("Dumping data to ", path)
	tmpMap := map[string]*TimeoutQueue{}
	p.Range(func(key, queue interface{}) bool {
		if !queue.(*TimeoutQueue).Empty() {
			tmpMap[key.(string)] = queue.(*TimeoutQueue)
		}
		return true
	})
	data, err := json.Marshal(tmpMap)
	if err != nil {
		log.Printf("Failed to marshal data from %s : %v\n", path, err)
		return
	}
	if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
		log.Printf("Failed to dump file to %s : %v\n", path, err)
	}
}