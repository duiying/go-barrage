package handlers

import (
	"encoding/json"
	"github.com/wuqinqiang/chitchat/models"
	"html/template"
	"io"
	"net/http"
)

func ChatRoom(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	}
	messages, err := models.Messages()
	if err != nil {
		danger(err.Error())
	} else {
		generateHTML(w, messages, "layout", "auth.navbar", "chat")
	}
}

//接口返回数据

func MessageAll(w http.ResponseWriter, r *http.Request) {
	messages, err := models.Messages()
	if err != nil {
		danger(err.Error())
	}
	b, err := json.Marshal(messages)
	if err != nil {
		danger(err.Error())
	}
	io.WriteString(w, string(b))
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			warning("cannot read form")
		}

		user, err := sess.User()
		if err != nil {
			warning("cannot find user")
		}
		ip := RemoteIP(r)
		content := template.HTMLEscapeString(r.PostFormValue("message"))

		if _, err := user.CreateMessage(ip, content, 1); err != nil {
			warning("send message error")
		}
		http.Redirect(w, r, "/chat", 302)
	}

}
