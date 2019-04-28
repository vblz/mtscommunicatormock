package mtsMock

import (
	"errors"
	"github.com/vblz/mtscommunicatormock/mtsWsdl"
	"github.com/vblz/mtscommunicatormock/store"
	"time"
)

func (m *mtsMock) ProcessSendMessage(request *mtsWsdl.SendMessage) (*mtsWsdl.SendMessageResponse, error) {
	if request == nil || !m.isAuthenticated(request.Login, request.Password) || m.naming != request.Naming {
		return nil, authError
	}

	id, err := m.store.AddOutgoingMessage(store.Message{Sent: time.Now(), Text: request.Message, Phone: request.Msid})
	if err != nil {
		return nil, errors.New("database error")
	}
	return &mtsWsdl.SendMessageResponse{SendMessageResult: id}, nil
}

func (m *mtsMock) ProcessGetMessagesStatus(request *mtsWsdl.GetMessagesStatus) (*mtsWsdl.GetMessagesStatusResponse, error) {
	if request == nil || !m.isAuthenticated(request.Login, request.Password) {
		return nil, authError
	}

	messages := make([]*mtsWsdl.MessageStatusWithID, len(request.MessageIDs.Long))
	for i, id := range request.MessageIDs.Long {
		deliveryStatus := mtsWsdl.DeliveryStatusNotSent
		deliveryDate := time.Now()
		if sent, err := m.store.GetOutgoingMessage(id); err == nil {
			deliveryStatus = mtsWsdl.DeliveryStatusDelivered
			deliveryDate = sent.Sent
		}
		deliverInfo := []*mtsWsdl.DeliveryInfo{{DeliveryDate: deliveryDate, DeliveryStatus: &deliveryStatus}}
		delivery := mtsWsdl.ArrayOfDeliveryInfo{DeliveryInfo: deliverInfo}
		messages[i] = &mtsWsdl.MessageStatusWithID{MessageID: id, Delivery: &delivery}
	}
	result := mtsWsdl.ArrayOfMessageStatusWithID{MessageStatusWithID: messages[:]}

	return &mtsWsdl.GetMessagesStatusResponse{GetMessagesStatusResult: &result}, nil
}
