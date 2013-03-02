package connectiontest

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/poll", poll)
	http.HandleFunc("/push", push)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, homepage)
}

const homepage = `
<DOCTYPE html>
<html>
	<body>
		<form method="POST" action="/push?rcpt=aaa">
			<textarea name="body">bla</textarea>
			<input type="submit">
		</form>
	</body>
</html>
`

var messages map[string] chan string = make(map[string] chan string)

func push(w http.ResponseWriter, req *http.Request) {
	rcpt := req.FormValue("rcpt")
	body := req.FormValue("body")

	if rcpt == "" || req.Method != "POST" {
		w.WriteHeader(400)
		return
	}

	ch := messages[rcpt]

	// new client?
	if ch == nil {
		ch = make (chan string)
		messages[rcpt] = ch
	}
	
	// TODO check whether the channel has space? probably don't want this to block
	ch <- string(body)

	fmt.Fprint(w, homepage)
}

func poll(w http.ResponseWriter, req *http.Request) {
	select {
		case <-time.After(15e9):
        	return
		case msg := <-messages["aaa"]:
			io.WriteString(w, msg)
	}
}
