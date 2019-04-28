package mtsMock

import (
	"errors"
	"github.com/vblz/mtscommunicatormock/mtsWsdl"
)

func (m *mtsMock) ProcessGetMessages(request *mtsWsdl.GetMessages) (*mtsWsdl.GetMessagesResponse, error) {
	if request == nil || !m.isAuthenticated(request.Login, request.Password) {
		return nil, authError
	}

	if request.MessageType == nil || *request.MessageType != mtsWsdl.RequestMessageTypeMO {
		return nil, errors.New("incorrect message type")
	}

	if request.SubscriberMsids != nil && len(request.SubscriberMsids.String) > 0 {
		return nil, errors.New("incorrect SubscriberMsids")
	}

	incomingMessages, err := m.store.GetIncomingMessages(request.DateFrom, request.DateTo)
	if err != nil {
		return nil, errors.New("database error")
	}
	messages := make([]*mtsWsdl.MessageInfo, len(incomingMessages))

	for i, v := range incomingMessages {
		messages[i] = &mtsWsdl.MessageInfo{
			MessageID:    v.Id,
			MessageText:  v.Text,
			SenderMsid:   v.Phone,
			CreationDate: v.Sent,
		}
	}

	result := mtsWsdl.ArrayOfMessageInfo{MessageInfo: messages}

	return &mtsWsdl.GetMessagesResponse{GetMessagesResult: &result}, nil
}