package inMemory

import (
	"errors"
	"github.com/vblz/mtscommunicatormock/store"
	"math"
	"math/rand"
	"sync"
	"time"
)

type inMemory struct {
	currentId      int64
	sentMessages   map[int64]store.Message
	sentMessagesIds []int64
	incomeMessages []store.Message
	currentIdMux   sync.Mutex
	// should use RWMutex, but who cares
	sentMessagesMux   sync.Mutex
	incomeMessagesMux   sync.Mutex
}

func (inMemoryStore *inMemory) AddIncomingMessage(m store.Message) error {
	m.Sent = time.Now()
	m.Id = inMemoryStore.nextId()
	inMemoryStore.incomeMessagesMux.Lock()
	defer inMemoryStore.incomeMessagesMux.Unlock()
	inMemoryStore.incomeMessages = append(inMemoryStore.incomeMessages, m)
	return nil
}

func (inMemoryStore *inMemory) GetOutgoingMessages(s store.SearchRequest) ([]store.Message, error) {
	inMemoryStore.sentMessagesMux.Lock()
	defer inMemoryStore.sentMessagesMux.Unlock()
	l := len(inMemoryStore.sentMessagesIds)
	startIndex := l-int(s.Limit)
	if startIndex < 0 {
		startIndex = 0
	}
	result := make([]store.Message, 0, s.Limit)
	for i := l - 1; i >= startIndex; i-- {
		result = append(result, inMemoryStore.sentMessages[inMemoryStore.sentMessagesIds[i]])
	}

	return result, nil
}

func (inMemoryStore *inMemory) GetIncomingMessages(from time.Time, to time.Time) ([]store.Message, error) {
	inMemoryStore.incomeMessagesMux.Lock()
	defer inMemoryStore.incomeMessagesMux.Unlock()
	l := len(inMemoryStore.incomeMessages)
	fromIndex := l
	toIndex := l

	for i := l - 1; i >= 0 && inMemoryStore.incomeMessages[i].Sent.Unix() > from.Unix(); i-- {
		fromIndex = i
		if toIndex == l && inMemoryStore.incomeMessages[i].Sent.Unix() <= to.Unix() {
			toIndex = i
		}
	}

	if toIndex == l {
		return []store.Message{}, nil
	}

	return inMemoryStore.incomeMessages[fromIndex : toIndex+1], nil
}

func (inMemoryStore *inMemory) AddOutgoingMessage(m store.Message) (int64, error) {
	m.Id = inMemoryStore.nextId()
	inMemoryStore.sentMessagesMux.Lock()
	defer inMemoryStore.sentMessagesMux.Unlock()
	inMemoryStore.sentMessages[m.Id] = m
	inMemoryStore.sentMessagesIds = append(inMemoryStore.sentMessagesIds, m.Id)
	return m.Id, nil
}

func (inMemoryStore *inMemory) GetOutgoingMessage(id int64) (store.Message, error) {
	if result, ok := inMemoryStore.sentMessages[id]; ok {
		return result, nil
	}

	return store.Message{}, errors.New("not found")
}

func NewInMemory() *inMemory {
	start := rand.NewSource(time.Now().UnixNano()).Int63()
	return &inMemory{
		currentId:    start,
		sentMessages: make(map[int64]store.Message)}
}

func (inMemoryStore *inMemory) nextId() int64 {
	inMemoryStore.currentIdMux.Lock()
	defer inMemoryStore.currentIdMux.Unlock()
	if inMemoryStore.currentId == math.MaxInt64 {
		inMemoryStore.currentId = 0
	}
	inMemoryStore.currentId += 1
	return inMemoryStore.currentId
}
