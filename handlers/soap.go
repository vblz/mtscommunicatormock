package handlers

import (
	"encoding/xml"
	"errors"
	"github.com/go-pkgz/lgr"
	"github.com/hooklift/gowsdl/soap"
	"github.com/vblz/mtscommunicatormock/mtsMock"
	"github.com/vblz/mtscommunicatormock/mtsWsdl"
	"github.com/vblz/mtscommunicatormock/store"
	"io/ioutil"
	"net/http"
	"strings"
)

const headerSoapaction = "Soapaction"
const prefixMtsCommunicator = "http://mcommunicator.ru/M2M/"

type handler struct {
	mts   mtsMock.MtsMock
	store store.Store
}

func NewHandler(m mtsMock.MtsMock, store store.Store) *handler {
	return &handler{mts: m, store: store}
}

func (s *handler) SoapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	soapAction := r.Header.Get(headerSoapaction)
	if soapAction[0] == '"' && soapAction[len(soapAction)-1] == '"' {
		soapAction = soapAction[1 : len(soapAction)-1]
	}

	if !strings.HasPrefix(soapAction, prefixMtsCommunicator) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	method := soapAction[len(prefixMtsCommunicator):]
	var request interface{}
	var process func(interface{}) (interface{}, error)
	switch method {
	case "SendMessage":
		request = new(mtsWsdl.SendMessage)
		process = func(r interface{}) (interface{}, error) { return s.mts.ProcessSendMessage(r.(*mtsWsdl.SendMessage)) }
	case "GetMessages":
		request = new(mtsWsdl.GetMessages)
		process = func(r interface{}) (interface{}, error) { return s.mts.ProcessGetMessages(r.(*mtsWsdl.GetMessages)) }
	case "GetMessagesStatus":
		request = new(mtsWsdl.GetMessagesStatus)
		process = func(r interface{}) (interface{}, error) { return s.mts.ProcessGetMessagesStatus(r.(*mtsWsdl.GetMessagesStatus)) }
	default:
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err := readSoapRequest(r, request)
	if err != nil {
		lgr.Printf("[INFO] error while deserialize request: %s", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	response, err := process(request)
	if err != nil {
		if s.mts.IsAuthError(err) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	writeSoapResponse(w, response)
}

func readSoapRequest(r *http.Request, result interface{}) error {
	respEnvelope := new(soap.SOAPEnvelope)
	respEnvelope.Body = soap.SOAPBody{Content: result}
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(rawBody) == 0 {
		return errors.New("пустое тело запроса")
	}

	err = xml.Unmarshal(rawBody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}

func writeSoapResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Add("Content-Type", "text/xml; charset=utf-8")
	envelope := soap.SOAPEnvelope{}

	envelope.Body.Content = response

	encoder := xml.NewEncoder(w)

	if err := encoder.Encode(envelope); err != nil {
		lgr.Printf("[ERROR] error while serialize response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := encoder.Flush(); err != nil {
		lgr.Printf("[ERROR] error while write response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
