package data

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/util"
	"sync"
	"time"
)

var (
	logger   = util.Logger()
	instance *EventDataRegistry
	once     sync.Once
)

type EventDataRegistry struct {
	m             sync.Mutex
	registryItems map[*websocket.Conn]*RegistryItem
}

// GetDataRegistry returns the singleton instance
func GetDataRegistry() *EventDataRegistry {
	once.Do(func() {
		instance = &EventDataRegistry{
			registryItems: make(map[*websocket.Conn]*RegistryItem),
		}
	})
	return instance
}

func (r *EventDataRegistry) MarkReady(
	client *websocket.Conn,
	startRequest *domain.CommandRequest,
	roomInfo *domain.RoomInfo) {
	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(roomInfo)
	logger.Infof("🚘MarkReady invoked connection addr:%s room:%s", client.RemoteAddr(), marshal)

	r.registryItems[client] = &RegistryItem{
		conn:          client,
		StartRequest:  *startRequest,
		RoomInfo:      *roomInfo,
		Statistics:    domain.InitStatisticStruct(),
		PlayDeque:     util.NewFixedSizeDeque(1024),
		RuleGroupList: make(map[domain.CounterType][]domain.Rule),
	}
}

func (r *EventDataRegistry) IsReady(client *websocket.Conn) bool {
	r.m.Lock()
	defer r.m.Unlock()
	_, ok := r.registryItems[client]
	return ok
}

func (r *EventDataRegistry) RemoveClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.registryItems, client)
	logger.Info("✂️[EventDataRegistry]RemoveClient invoked connection addr:", client.RemoteAddr())
}

// StartPushPlayMessage 循环检查play队列，出队并推送
func (r *EventDataRegistry) StartPushPlayMessage() {
	round := 0
	for {
		if round%constant.LogRound == 0 {
			logger.Infof("[EventDataRegistry]start StartPushPlayMessage[🎬] checking....ready clients size: %d", len(r.registryItems))
			round = 0
		}

		// Check each connected connection
		for client, registryItem := range r.registryItems {
			r.pushUserAction(client, registryItem)
		}

		// Sleep for PlayDequeuePushInterval(=1 by default) seconds
		time.Sleep(constant.PlayDequeuePushInterval * time.Second)
		round++

	}

}

// pop action queue and query NLP to try get a response
func (r *EventDataRegistry) pushUserAction(client *websocket.Conn, registryItem *RegistryItem) {
	for !registryItem.PlayDeque.IsEmpty() {
		playMessage := registryItem.PlayDeque.PopFront()
		if playMessage.Action != domain.ON_COMMENT {
			continue
		}
		if !constant.PlayUserAction {
			continue
		}

		item := r.registryItems[client]
		queryResponse := item.RequestNlp(playMessage.Username, playMessage.Content)
		if queryResponse == nil || queryResponse.Meta.Catchall {
			logger.Infof("[EventDataRegistry]PushPlayMessage[🤐]skip for CHATCHALL to  query: %s connection: %s", playMessage.Content, client.RemoteAddr())
			continue
		}

		message := queryResponse.ToPlayMessage()
		marshal, _ := json.Marshal(playMessage)
		logger.Infof("[EventDataRegistry]PushPlayMessage[🎬⚙️] chat message: %s connection: %s", marshal, client.RemoteAddr())
		r.WriteResponse(client, message)
	}
}

func (r *EventDataRegistry) LoadRule(traceId string, client *websocket.Conn) constant.ResponseStatus {
	item := r.registryItems[client]
	if item == nil {
		return constant.CLIENT_NOT_READY
	}
	return item.LoadRule(traceId)
}

func (r *EventDataRegistry) WriteResponse(client *websocket.Conn, commandResponse *domain.CommandResponse) error {
	item := r.registryItems[client]
	if item == nil {
		//commandResponse.ResponseStatus = constant.CLIENT_NOT_READY
		marshal, _ := json.Marshal(commandResponse)
		client.WriteMessage(websocket.TextMessage, marshal)
		return errors.New(constant.CLIENT_NOT_READY.Message)
	}
	err := item.WriteResponse(commandResponse)
	if err != nil {
		logger.Errorf("[dataRegistry]WriteResponse fail with err:%s", err)
	}
	return nil
}

func (r *EventDataRegistry) UpdateStatistics(conn *websocket.Conn, counterType domain.CounterType, counter *domain.StatisticCounter) error {
	item := r.registryItems[conn]
	if item == nil {
		return errors.New(constant.CLIENT_NOT_READY.Message)
	}

	item.UpdateStatistics(counterType, counter)
	return nil
}

func (r *EventDataRegistry) EnqueueAction(conn *websocket.Conn, event domain.UserActionEvent) error {
	item := r.registryItems[conn]
	if item == nil {
		return errors.New(constant.CLIENT_NOT_READY.Message)
	}

	item.EnqueueAction(event)
	return nil
}
