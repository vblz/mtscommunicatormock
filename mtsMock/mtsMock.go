package mtsMock

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/vblz/mtscommunicatormock/mtsWsdl"
	"github.com/vblz/mtscommunicatormock/store"
	"strings"
)

type MtsMock interface {
	ProcessSendMessage(request *mtsWsdl.SendMessage) (*mtsWsdl.SendMessageResponse, error)
	ProcessGetMessagesStatus(request *mtsWsdl.GetMessagesStatus) (*mtsWsdl.GetMessagesStatusResponse, error)
	ProcessGetMessages(request *mtsWsdl.GetMessages) (*mtsWsdl.GetMessagesResponse, error)
	IsAuthError(err error) bool
}

type mtsMock struct {
	store        store.Store
	login        string
	passwordHash string
	naming       string
}

var authError =  errors.New("auth error")

func NewMtsMock(login, password, naming string, store store.Store) MtsMock {
	return &mtsMock{login: login, passwordHash: getMD5Hash(password), store: store, naming: naming}
}

func (m *mtsMock) isAuthenticated(login, passwordHash string) bool {
	// it's unsafe compare in real code due timing attack, but it should not use with real claims. Never
	return strings.EqualFold(m.login, login) && strings.EqualFold(m.passwordHash, passwordHash)
}

func (m *mtsMock) IsAuthError(err error) bool {
	return err == authError
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
