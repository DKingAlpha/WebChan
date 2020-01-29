package tpl

import (
	"github.com/DKingCN/WebChan/internal/shared_vars"
	"html/template"
	"log"
)

var TplAdmin *template.Template = nil

func init() {
	t, err := template.New("admin").Parse(tplAdmin)
	if err != nil {
		log.Fatalf("Failed to init tpl: %v\n", err)
	}
	TplAdmin = t
}

const tplAdmin string = `<html>
<head>
<script>
function delete_channel(channel){
	if(!confirm('确认删除?'))return;
    var xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function () {
		if (this.readyState != 4) return;
		if (this.status == 200) {
			location.reload();
		}
	};
	xhr.open("DELETE", "/"+channel+"/?admin="+"`+ shared_vars.AdminPassword +`", false);
	xhr.send();
}
</script>
</head>

<body>
<table border=1 cellpadding=8 cellspacing=0>
<tr>
	<th>Channel</th>
	<th>Key</th>
	<th>Perm</th>
</tr>
{{range $channelId, $rtq := .}}
<tr>
	<td><a onclick="delete_channel('{{$channelId}}')">[x]</a> <a href="/{{$channelId}}?key={{$rtq.Key}}">{{$channelId}}</a></td>
	<td>{{$rtq.Key}}</td>
	<td>{{$rtq.Perm.String}}</td>
</tr>
{{end}}
</table>
</body>
</html>
`


