package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// 7 days activity
var activityLog *ActivityLog = nil



func statusHandler(w http.ResponseWriter, r *http.Request) {
	tpl := `<html>
<table border=1 cellpadding=8 cellspacing=0>
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

<br><br>
<pre>
# Usage
curl -X POST  qaq.link/chan1/msg1
curl -X POST  qaq.link/chan1  --data-binary msg2
curl -X POST  qaq.link/chan1/msg2?key=PASSWORD
</pre>
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
