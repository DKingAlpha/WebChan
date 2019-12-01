package tpl

import (
	"html/template"
	"log"
)

var TplIndex *template.Template = nil

func init() {
	t, err := template.New("index").Parse(tplIndex)
	if err != nil {
		log.Fatalf("Failed to init tpl: %v\n", err)
	}
	TplIndex = t
}


const tplIndex string = `<html>
<body>
<table border=1 cellpadding=8 cellspacing=0>
<tr>
	<th>Channel</th>
	<th>Count</th>
	<th>Last Activity</th>
</tr>
{{range .}}
<tr>
	<td><a href="/{{.Chan}}">{{.Chan}}</a></td>
	<td>{{.Act.Count}}</td>
	<td>{{.Act.GetTime}}</td>
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
</body>
</html>`