package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// 7 days activity
var activityLog *ActivityLog = nil

func adminHandler(w http.ResponseWriter, r *http.Request) {
	tmpMap := map[string]*RTQ{}
	queues.Range(func(key, queue interface{}) bool {
		if !queue.(*RTQ).Empty() {
			tmpMap[key.(string)] = queue.(*RTQ)
		}
		return true
	})
	tpl := `<html>
<table border=1 cellpadding=8 cellspacing=0>
<tr>
	<th>Channel</th>
	<th>Key</th>
	<th>Perm</th>
</tr>
{{range $channelId, $rtq := .}}
<tr>
	<th><a href="/{{$channelId}}?key={{$rtq.Key}}">{{$channelId}}</a></th>
	<th>{{$rtq.Key}}</th>
	<th>{{$rtq.Perm.String}}</th>
</tr>
{{end}}
</table>
</html>
`
	tmpl, err := template.New("admin").Parse(tpl)
	if err != nil {
		log.Fatalf("Failed to init tpl for status: %v\n", err)
	}
	err = tmpl.Execute(w, tmpMap)
	if err != nil {
		_,_ = fmt.Fprintf(w, "Failed to get status: %v\n", err)
	}
}

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
### What is this
A data exchanging server, based on HTTP requests.

Data are organized by "channel", with "timeout", "capability" and "access control" enforced.

### Rules

##### Basic
1. Basic request pattern:   POST /channel/msg1  msg2(in body)    |   GET /channel   |   DELETE /channel
2. POST writes messages to channel, GET reads messages from channel, DELETE delete a channel
3. POST: both msg1&msg2 is optional. If neither of them exists, it means creating a empty channel

##### Access control
4. POST with params "key" and "perm", to set a key to channel, and/or set permissions for others(who POST/GET/DELETE without a key). A key makes channel private
5. "perm" consist of "r w d" which stands for read/write/delete. for example, perm=r means "chmod o+r channel"
6. append params to URL in form of "?key=passwd&perm=rwd". Or ignore any of them to use default value
7. GET/DELETE with key if its a private channel and you dont have related permission

##### Things you might want to know
8. Anyone could update key and perm of a public channel
9. For a public  channel, default value: key is empty, perm=rwd
10. For a private channel, default value: perm is none
11. msg_timeout=len(channel_id) days
12. Timeout msg will be removed from channel. Channel without msgs will be released(DELETE) soon
13. You may feel like to use PUT instead of POST, or OPTIONS while updating key/perm. They are just alias of POST here.

#### Example
# create a new public channel
curl -X POST  qaq.link/chan1
curl -X POST  qaq.link/chan1/

# post msg to a channel(create if not exists)
curl -X POST  qaq.link/chan1/msg1
curl -X POST  qaq.link/chan1  --data-binary msg2
curl -X POST  qaq.link/chan1/msg1  --data-binary msg2
curl -X POST  qaq.link/chan1/msg2?key=PASSWORD
curl -X POST  qaq.link/chan1/msg2?key=PASSWORD&perm=rwd

# get data
curl -X GET qaq.link/chan1/msg1
curl qaq.link/chan1/msg1
curl qaq.link/chan1/msg2?key=PASSWORD
</pre>
</html>`
	tmpl, err := template.New("status").Parse(tpl)
	if err != nil {
		log.Fatalf("Failed to init tpl for status: %v\n", err)
	}
	chanActs := activityLog.Rank()
	err = tmpl.Execute(w, *chanActs)
	if err != nil {
		_,_ = fmt.Fprintf(w, "Failed to get status: %v\n", err)
	}
}
