package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// 7 days activity
var activityLog = ActivityLog{
	Cap:     100,
	Timeout: 1*24*60*60,
	Acts:    map[string]*Activity{},
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	tpl := `<html>
<table>
<tr>
	<th>Channel</th>
	<th>Count</th>
	<th>Last Activity</th>
</tr>
{{range .}}
<tr>
	<th><a href="/{{.Chan}}">{{.Chan}}</a></th>
	<th>{{.Act.Count}}</th>
	<th>{{.Act.GetTime}}</th>
</tr>
{{end}}

</table>
</html>`
	tmpl, err := template.New("status").Parse(tpl)
	if err != nil {
		log.Fatal("Failed to init tpl for status")
	}
	chanActs := activityLog.Rank()
	err = tmpl.Execute(w, *chanActs)
	if err != nil {
		_,_ = fmt.Fprintf(w, "Failed to get status: %v\n", err)
	}
}
