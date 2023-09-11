package data

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/util"
	"sync"
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

func (r *EventDataRegistry) Size() int {
	return len(r.registryItems)
}

func (r *EventDataRegistry) MarkReady(
	client *websocket.Conn,
	startRequest *domain.CommandRequest,
	roomInfo *domain.RoomInfo,
	stopChan chan struct{}) {

	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(roomInfo)
	logger.Infof("🚘MarkReady invoked connection addr:%s room:%s", client.RemoteAddr(), marshal)

	registryItem := &RegistryItem{
		conn:          client,
		StartRequest:  *startRequest,
		RoomInfo:      *roomInfo,
		Statistics:    domain.InitStatisticStruct(),
		PlayDeque:     util.NewFixedSizeDeque(1024),
		RuleGroupList: make(map[domain.CounterType][]domain.Rule),
		stopChan:      stopChan,
	}
	r.registryItems[client] = registryItem
	go registryItem.StartPushPlayMessage()
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
	addr := ""
	if conn != nil {
		addr = conn.LocalAddr().String()
	}
	logger.Infof("UpdateStatistics LocalAddr:%s CounterType:%s counter:%s", addr, counterType, counter)
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
