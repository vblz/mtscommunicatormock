package handlers

import (
	"github.com/go-pkgz/lgr"
	"github.com/vblz/mtscommunicatormock/store"
	"html/template"
	"net/http"
)

const sendLocation = "static/send.html"
const listLocation = "static/list.html"

var sendTemplate = template.Must(template.ParseFiles(sendLocation))
var listTemplate = template.Must(template.ParseFiles(listLocation))

func (s *handler) SendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := sendTemplate.Execute(w, nil)
		if err != nil {
			lgr.Printf("[ERROR] send template execute error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		lgr.Printf("[INFO] send parse form error: %s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	phones := r.Form["phone"]
	texts := r.Form["text"]

	if len(phones) != 1 || len(texts) != 1 {
		lgr.Printf("[INFO] incorrect user input")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	phone := phones[0]
	text := texts[0]

	if phone == "" || text == "" {
		lgr.Printf("[INFO] empty input")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = s.store.AddIncomingMessage(store.Message{Phone: phone, Text: text})
	if err != nil {
		lgr.Printf("[ERROR] error while AddIncomingMessage: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.URL.String(), 303)
}

func (s *handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		lgr.Printf("[INFO] list handler not get method")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	searchParams := store.SearchRequest{Limit:10}

	messages, err := s.store.GetOutgoingMessages(searchParams)
	if err != nil {
		lgr.Printf("[ERROR] error while GetOutgoingMessages: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = listTemplate.Execute(w, messages)
	if err != nil {
		lgr.Printf("[ERROR] can't execute template list: %s", err)
	}
}
