package registry

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"live-room-crawler/platform"
	"live-room-crawler/util"
	"sync"
	"time"
)

var (
	logger   = util.Logger()
	instance *LocalClientRegistry
	once     sync.Once
)

type LocalClientRegistry struct {
	m                     sync.Mutex
	clients               map[*websocket.Conn]int
	HeartbeatLostRegistry map[*websocket.Conn]int
	registryClients       map[*websocket.Conn]ClientRegistryItem
}

type ClientRegistryItem struct {
	lostHeatBeatStamp         int64
	StartRequest              common.CommandRequest
	RoomInfo                  common.RoomInfo
	PlayDeque                 common.PlayDeque
	Statistics                common.LiveStatisticsStruct
	PlatformConnectorStrategy platform.IPlatformConnectorStrategy
}

// GetInstance returns the singleton instance
func GetInstance() *LocalClientRegistry {
	once.Do(func() {
		instance = &LocalClientRegistry{
			clients:         make(map[*websocket.Conn]int),
			registryClients: make(map[*websocket.Conn]ClientRegistryItem),
		}
	})
	return instance
}

func (r *LocalClientRegistry) AddClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	logger.Info("AddClient invoked client addr:", client.RemoteAddr())

	r.clients[client] = 0
	r.HeartbeatLostRegistry[client] = 0
}

func (r *LocalClientRegistry) MarkReady(
	client *websocket.Conn,
	startRequest common.CommandRequest,
	connectorStrategy platform.IPlatformConnectorStrategy) {
	r.m.Lock()
	defer r.m.Unlock()
	info := connectorStrategy.GetRoomInfo()
	marshal, _ := json.Marshal(info)
	logger.Infof("markReady invoked client addr:%s room:%s", client.RemoteAddr(), marshal)

	r.registryClients[client] = ClientRegistryItem{
		StartRequest:              startRequest,
		RoomInfo:                  *info,
		Statistics:                common.InitStatisticStruct(),
		PlayDeque:                 *common.NewFixedSizeDeque(1024),
		PlatformConnectorStrategy: connectorStrategy,
	}
}

func (r *LocalClientRegistry) RemoveClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	logger.Info("RemoveClient invoked client addr:", client.RemoteAddr())
	delete(r.clients, client)
	delete(r.registryClients, client)
}

func (r *LocalClientRegistry) UpdateStatistics(client *websocket.Conn, statistics common.LiveStatisticsStruct) {
	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(statistics)
	logger.Infof("UpdateStatistics invoked client addr:%s statistics:%s", client.RemoteAddr(), marshal)

	item := r.registryClients[client]
	item.Statistics.Add(statistics)
}

func (r *LocalClientRegistry) EnqueueAction(client *websocket.Conn, actionEvent common.UserActionEvent) {
	r.m.Lock()
	defer r.m.Unlock()
	marshal, _ := json.Marshal(actionEvent)
	logger.Infof("EnqueueAction invoked client addr:%s event:%s", client.RemoteAddr(), marshal)

	item := r.registryClients[client]
	item.PlayDeque.PushBack(actionEvent)
}

func (r *LocalClientRegistry) DequeueAction(client *websocket.Conn) *common.UserActionEvent {
	r.m.Lock()
	defer r.m.Unlock()

	item := r.registryClients[client]
	front := item.PlayDeque.PopFront()

	marshal, _ := json.Marshal(front)
	logger.Infof("DequeueAction invoked client addr:%s event:%s", client.RemoteAddr(), marshal)
	return front
}

// StartHeartbeatsCheck 检查心跳
func (r *LocalClientRegistry) StartHeartbeatsCheck() {
	for {
		logger.Info("[local_server]start StartHeartbeatsCheck[❤️]....")
		// Sleep for 10 seconds
		time.Sleep(30 * time.Second)

		// Check each connected client
		for client, missedHeartbeats := range r.HeartbeatLostRegistry {
			if missedHeartbeats > 3 {
				// If the client has missed more than 3 heartbeat_losts, remove it
				logger.Infoln("[local_server]remove client for heartbeat timeout[🕔💔] address:", client.RemoteAddr())
				r.RemoveClient(client)
			} else {
				// Otherwise, send a heartbeat and increment the missed heartbeat count
				err := client.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					logger.Error("[local_server]fail to write pong to client", err)
					r.RemoveClient(client)
				} else {
					r.HeartbeatLostRegistry[client]++
				}
			}
		}
	}
}

// StartPushPlayMessage 循环检查play队列，出队并推送
func (r *LocalClientRegistry) StartPushPlayMessage() {
	for {
		logger.Info("[local_server]start StartPushPlayMessage[🎬]....")
		// Sleep for 10 seconds
		time.Sleep(30 * time.Second)

		// Check each connected client
		for client, registryItem := range r.registryClients {
			playMessage := registryItem.PlayDeque.PopFront()
			if playMessage != nil {
				marshal, _ := json.Marshal(playMessage)
				logger.Infof("[local_server]PushPlayMessage[🎬⚙️] message: %s client: %s", marshal, client.RemoteAddr())
				client.WriteMessage(websocket.TextMessage, marshal)
			}

		}

	}

}
