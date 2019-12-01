package tpl

/*
var TplTool *template.Template = nil

func init() {
	t, err := template.New("tool").Parse(tplTool)
	if err != nil {
		log.Fatalf("Failed to init tpl: %v\n", err)
	}
	TplTool = t
}
*/


const TplTool string = `
<html>
<head>
<script>
function send_to_channel(){
	let channel=document.getElementById("channel").value;
	let data=document.getElementById("data").value;
	let key=document.getElementById("key").value;
	let perm="";
	if(document.getElementById("perm_r").checked)perm+="r";
	if(document.getElementById("perm_w").checked)perm+="w";
	if(document.getElementById("perm_d").checked)perm+="d";
	console.log(perm);

    var xhr = new XMLHttpRequest();
	let channelUrl = "/" + channel + "?key=" + key + "&perm=" + perm;
	xhr.onreadystatechange = function () {
		if (this.readyState != 4) return;
		if (this.status == 200) {
			location=channelUrl;
		}
	};
	xhr.open("POST", channelUrl);
	xhr.send(data);
}
</script>
</head>

<body>
<form>
  <p>Channel: <input type="text" id="channel" /></p>
  <p>Key: <input type="text" id="key" /></p>
  <p>Perm: R<input type="checkbox" id="perm_r" checked="checked"/> W<input type="checkbox" id="perm_w" /> D<input type="checkbox" id="perm_d" /></p>
  <textarea id="data"></textarea><br>
  <input type="submit" onclick="send_to_channel()"/>
</form>

</body>
</html>
`