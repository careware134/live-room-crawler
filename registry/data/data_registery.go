package data

import (
	"encoding/json"
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
	registryItems map[*websocket.Conn]*DataRegistryItem
}

// GetDataRegistry returns the singleton instance
func GetDataRegistry() *EventDataRegistry {
	once.Do(func() {
		instance = &EventDataRegistry{
			registryItems: make(map[*websocket.Conn]*DataRegistryItem),
		}
	})
	return instance
}

type DataRegistryItem struct {
	lostHeatBeatStamp int64
	StartRequest      common.CommandRequest
	RoomInfo          common.RoomInfo
	PlayDeque         *common.PlayDeque
	Statistics        common.LiveStatisticsStruct
	//Client            server.LocalClient
	//PlatformConnectorStrategy platform.IPlatformConnectorStrategy
}

func (r *EventDataRegistry) MarkReady(
	client *websocket.Conn,
	startRequest *common.CommandRequest,
	roomInfo *common.RoomInfo) {
	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(roomInfo)
	logger.Infof("🚘MarkReady invoked connection addr:%s room:%s", client.RemoteAddr(), marshal)

	r.registryItems[client] = &DataRegistryItem{
		StartRequest: *startRequest,
		RoomInfo:     *roomInfo,
		Statistics:   common.InitStatisticStruct(),
		PlayDeque:    common.NewFixedSizeDeque(1024),
	}
}

func (r *EventDataRegistry) RemoveClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	delete(r.registryItems, client)
	logger.Info("✂️[EventDataRegistry]RemoveClient invoked connection addr:", client.RemoteAddr())
}

func (r *EventDataRegistry) UpdateStatistics(client *websocket.Conn,
	counterType common.CounterType,
	counter common.StatisticCounter) {
	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(counter)
	logger.Infof("UpdateStatistics for connection addr:%s type:%s counter:%s ", client.RemoteAddr(), counterType, marshal)

	item := r.registryItems[client]
	if item != nil {
		item.Statistics.AddCounter(counterType, counter)
	}
}

func (r *EventDataRegistry) EnqueueAction(client *websocket.Conn, actionEvent common.UserActionEvent) {
	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(actionEvent)
	logger.Infof("EnqueueAction invoked connection addr:%s event:%s", client.RemoteAddr(), marshal)

	item := r.registryItems[client]
	item.PlayDeque.PushBack(actionEvent)
}

func (r *EventDataRegistry) DequeueAction(client *websocket.Conn) *common.UserActionEvent {
	r.m.Lock()
	defer r.m.Unlock()

	item := r.registryItems[client]
	front := item.PlayDeque.PopFront()

	marshal, _ := json.Marshal(front)
	logger.Infof("DequeueAction invoked connection addr:%s event:%s", client.RemoteAddr(), marshal)
	return front
}

// StartPushPlayMessage 循环检查play队列，出队并推送
func (r *EventDataRegistry) StartPushPlayMessage() {
	for {
		logger.Infof("[EventDataRegistry]start StartPushPlayMessage[🎬] checking....ready clients size: %d", len(r.registryItems))

		// Check each connected connection
		for client, registryItem := range r.registryItems {
			for !registryItem.PlayDeque.IsEmpty() {
				playMessage := registryItem.PlayDeque.PopFront()
				if playMessage != nil {
					playMessage := playMessage.ToPlayMessage()
					marshal, _ := json.Marshal(playMessage)
					logger.Infof("[EventDataRegistry]PushPlayMessage[🎬⚙️] message: %s connection: %s", marshal, client.RemoteAddr())
					client.WriteMessage(websocket.TextMessage, marshal)
				}
			}
		}

		// Sleep for 10 seconds
		time.Sleep(constant.PlayDequeuePushInterval * time.Second)

	}

}
