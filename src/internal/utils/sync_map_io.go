package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

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
