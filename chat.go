package main

import (
	"container/list"
	"fmt"
	"go.net/websocket"
	"io"
	"net/http"
)

var connid int
var conns *list.List

func Chat(ws *websocket.Conn) {
	defer ws.Close()
	item := conns.PushBack(ws)
	defer conns.Remove(item)
	var err error
	for {
		var data string
		if err = websocket.Message.Receive(ws, &data); err != nil {
			fmt.Printf("disconnected\n")
			break
		}
		SendMessage(item, fmt.Sprintf("%s", data))
	}
}
func SendMessage(self *list.Element, data string) {
	//for _, item := range conns {
	for item := conns.Front(); item != nil; item = item.Next() {
		ws, ok := item.Value.(*websocket.Conn)
		if !ok {
			panic("item not *websocket.Conn")
		}
		if item == self {
			continue
		}
		io.WriteString(ws, data)
	}
}

// 客户端默认显示页面
func Client(w http.ResponseWriter, r *http.Request) {
	html := `<!doctype html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>golang websocket chatroom</title>
    <script language="javascript"type="text/javascript">  
        var sock=null; 
        var wsuri ="ws://192.168.1.163:7878/chat"; //这里的IP如果是局域测试的话，需要换成自己的
        window.onload = function(){
            console.log("onload");
            sock = new WebSocket(wsuri);
            sock.onopen=function(e){
                console.log("connected to "+wsuri);
            }
            sock.onclose=function(e){
                console.log("connection closed (" + e.code + ")");
            }
            sock.onmessage=function(e){
                console.log("message received: " + e.data);
                document.getElementById("list").innerHTML += e.data;
                var div = document.getElementById("list")
                div.scrollTop = div.scrollHeight; 
            }
        }
        function send () {
            var who = document.getElementById('who').value;
            if (who.length == 0){
                document.getElementById('who').focus();
                return
            }
            var msg = document.getElementById('msg').value;
            if (msg.length == 0){
                document.getElementById('msg').focus();
                return
            }
            var data = who + "  say:  " + msg + "<br/><br/>";
            document.getElementById('msg').value="";
            document.getElementById('who').value="";
            var div = document.getElementById("list")
            div.innerHTML += data;
            div.scrollTop = div.scrollHeight; 
            sock.send(data);
        }        
    </script>
</head>
<body>
    <div id="list" style="height: 300px;overflow-y: scroll;border: 1px solid #CCC;">
    </div>
    <div>
        who are you 
        <input type="text" id="who" size="60" />
        your message
        <input type="text" id="msg" size="60" />
        <button onclick="send()">send</button>
    </div>
</body>
</html>`
	io.WriteString(w, html)
}
func main() {
	conns = list.New()
	http.Handle("/chat", websocket.Handler(Chat))
	http.HandleFunc("/", Client)
	err := http.ListenAndServe(":7878", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
