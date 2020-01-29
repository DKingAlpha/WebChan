package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

func LoadSyncMap(path string) *sync.Map {
	log.Println("Loading map from", path)
	f := sync.Map{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load map from %s: %v\n", path, err)
		return &f
	}
	tmpMap := map[string]*RTQ{}
	if err := json.Unmarshal(data, &tmpMap); err != nil {
		log.Printf("Failed to unmarshal map: %v\n", err)
		return &f
	}
	for key, value := range tmpMap {
		f.Store(key, value)
	}
	log.Printf("Loaded %s\n", path)
	return &f
}

func DumpSyncMap(p *sync.Map, path string) {
	log.Println("Dumping map to", path)
	tmpMap := map[string]*RTQ{}
	p.Range(func(key, queue interface{}) bool {
		queue.(*RTQ).CleanTimeout()
		if !queue.(*RTQ).Empty() {
			tmpMap[key.(string)] = queue.(*RTQ)
		}
		return true
	})
	data, err := json.Marshal(tmpMap)
	if err != nil {
		log.Printf("Failed to marshal map: %v\n", err)
		return
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		log.Printf("Failed to dump map to %s: %v\n", path, err)
	}
}
