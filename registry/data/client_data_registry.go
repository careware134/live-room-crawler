package data

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"live-room-crawler/constant"
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
	startRequest *common.CommandRequest,
	roomInfo *common.RoomInfo) {
	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(roomInfo)
	logger.Infof("🚘MarkReady invoked connection addr:%s room:%s", client.RemoteAddr(), marshal)

	r.registryItems[client] = &RegistryItem{
		conn:          client,
		StartRequest:  *startRequest,
		RoomInfo:      *roomInfo,
		Statistics:    common.InitStatisticStruct(),
		PlayDeque:     common.NewFixedSizeDeque(1024),
		RuleGroupList: make(map[common.CounterType][]common.Rule),
	}
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

		// Check each connected connection
		//for client, registryItem := range r.registryItems {
		//	r.pushStatisticRulePlay(client, registryItem)
		//}

		// Sleep for 10 seconds
		time.Sleep(constant.PlayDequeuePushInterval * time.Second)
		round++

	}

}

func (r *EventDataRegistry) pushUserAction(client *websocket.Conn, registryItem *RegistryItem) {
	for !registryItem.PlayDeque.IsEmpty() {
		playMessage := registryItem.PlayDeque.PopFront()
		if constant.PlayUserAction && playMessage != nil {
			playMessage := playMessage.ToPlayMessage()
			marshal, _ := json.Marshal(playMessage)
			logger.Infof("[EventDataRegistry]PushPlayMessage[🎬⚙️] message: %s connection: %s", marshal, client.RemoteAddr())
			r.WriteResponse(client, playMessage)
		}
	}
}

func (r *EventDataRegistry) LoadRule(client *websocket.Conn) constant.ResponseStatus {
	item := r.registryItems[client]
	if item == nil {
		return constant.CLIENT_NOT_READY
	}
	return item.LoadRule()
}

func (r *EventDataRegistry) WriteResponse(client *websocket.Conn, message *common.CommandResponse) error {
	item := r.registryItems[client]
	if item == nil {
		message.ResponseStatus = constant.CLIENT_NOT_READY
		marshal, _ := json.Marshal(message)
		client.WriteMessage(websocket.TextMessage, marshal)
		return errors.New(constant.CLIENT_NOT_READY.Message)
	}
	item.WriteResponse(message)
	return nil
}

func (r *EventDataRegistry) UpdateStatistics(conn *websocket.Conn, counterType common.CounterType, counter *common.StatisticCounter) error {
	item := r.registryItems[conn]
	if item == nil {
		return errors.New(constant.CLIENT_NOT_READY.Message)
	}

	item.UpdateStatistics(counterType, counter)
	return nil
}

func (r *EventDataRegistry) EnqueueAction(conn *websocket.Conn, event common.UserActionEvent) error {
	item := r.registryItems[conn]
	if item == nil {
		return errors.New(constant.CLIENT_NOT_READY.Message)
	}

	item.EnqueueAction(event)
	return nil
}
