package main

import (
	"github.com/vblz/mtscommunicatormock/handlers"
	"github.com/vblz/mtscommunicatormock/mtsMock"
	"github.com/vblz/mtscommunicatormock/store/inMemory"
	"log"
	"net/http"
)

func main() {
	store := inMemory.NewInMemory()
	mts := mtsMock.NewMtsMock("my_login", "my_password", "my_naming", store)
	handler := handlers.NewHandler(mts, store)
	http.HandleFunc("/test.svc", handler.SoapHandler)
	http.HandleFunc("/ui/send", handler.SendHandler)
	http.HandleFunc("/ui/list", handler.ListHandler)
	log.Printf("Starting")
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
