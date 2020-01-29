package main

import (
	"encoding/json"
	"fmt"
	"github.com/DKingCN/WebChan/internal/tpl"
	"net/http"
)

var activityLog *ActivityLog = nil

func adminHandler(w http.ResponseWriter, r *http.Request) {
	tmpMap := map[string]*RTQ{}
	queues.Range(func(key, queue interface{}) bool {
		if !queue.(*RTQ).Empty() {
			tmpMap[key.(string)] = queue.(*RTQ)
		}
		return true
	})
	err := tpl.TplAdmin.Execute(w, tmpMap)
	if err != nil {
		_,_ = fmt.Fprintf(w, "Failed to get status: %v\n", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	chanActs := activityLog.Rank()
	err := tpl.TplIndex.Execute(w, *chanActs)
	if err != nil {
		_,_ = fmt.Fprintf(w, "Failed to get status: %v\n", err)
	}
}

func toolHandler(w http.ResponseWriter, r *http.Request) {
	// p := strings.Split(r.URL.Path, "/")[1:]
	_, _ = fmt.Fprintf(w, tpl.TplTool)
}

func recentHandler(w http.ResponseWriter, r *http.Request) {
	chanActs := activityLog.Rank()
	jsondata, err := json.Marshal(chanActs)
	if err == nil {
		_, _ = fmt.Fprintf(w, string(jsondata))
	}
}