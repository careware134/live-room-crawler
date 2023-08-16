package connection

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/registry/data"
	"live-room-crawler/util"
	"sync"
	"time"
)

var (
	logger   = util.Logger()
	instance *ClientConnectionRegistry
	once     sync.Once
)

type ClientConnectionRegistry struct {
	m                     sync.Mutex
	clients               map[*websocket.Conn]int
	heartbeatLostRegistry map[*websocket.Conn]int
	readyLocalClients     map[*websocket.Conn]*LocalClient
}

// GetClientRegistry returns the singleton instance
func GetClientRegistry() *ClientConnectionRegistry {
	once.Do(func() {
		instance = &ClientConnectionRegistry{
			clients:               make(map[*websocket.Conn]int),
			heartbeatLostRegistry: make(map[*websocket.Conn]int),
			readyLocalClients:     make(map[*websocket.Conn]*LocalClient),
		}
	})
	return instance
}

func (r *ClientConnectionRegistry) AddClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	logger.Info("AddClient invoked connection addr:", client.RemoteAddr())

	r.clients[client] = 0
	r.heartbeatLostRegistry[client] = 0
}

func (r *ClientConnectionRegistry) MarkReady(
	client *websocket.Conn,
	startRequest *domain.CommandRequest,
	localClient *LocalClient) {
	r.m.Lock()
	defer r.m.Unlock()
	roomInfo := (*localClient.Connector).GetRoomInfo()
	marshal, _ := json.Marshal(roomInfo)
	logger.Infof("🚘MarkReady invoked connection addr:%s room:%s", client.RemoteAddr(), marshal)

	r.clients[client] = 0
	r.readyLocalClients[client] = localClient
	r.heartbeatLostRegistry[client] = 0

	dataRegistry := data.GetDataRegistry()
	dataRegistry.MarkReady(client, startRequest, roomInfo)
}

func (r *ClientConnectionRegistry) UpdateHeartBeat(client *websocket.Conn) {
	logger.Info("UpdateHeartBeat for connection addr:", client.RemoteAddr())
	r.heartbeatLostRegistry[client] = 0
}

func (r *ClientConnectionRegistry) RemoveClient(client *websocket.Conn, tryRevoke bool) {
	r.m.Lock()
	defer r.m.Unlock()
	logger.Info("✂️[ClientConnectionRegistry]RemoveClient invoked connection addr:", client.RemoteAddr())
	delete(r.clients, client)
	localClient := r.readyLocalClients[client]
	if tryRevoke && localClient != nil {
		localClient.TryRevoke()
	}

	delete(r.readyLocalClients, client)
	delete(r.heartbeatLostRegistry, client)

	dataRegistry := data.GetDataRegistry()
	dataRegistry.RemoveClient(client)
}

// StartHeartbeatsCheck 检查心跳
func (r *ClientConnectionRegistry) StartHeartbeatsCheck() {
	for {
		logger.Debug("[ConnectionRegistry]start StartHeartbeatsCheck[❤️]....")
		// Sleep for 10 seconds
		time.Sleep(constant.HeartbeatCheckInterval * time.Second)

		// Check stop connector
		r.evictConnectorLost()

		// Check each connected connection
		r.evictHeartBeatLost()

	}
}

func (r *ClientConnectionRegistry) evictHeartBeatLost() {
	for client, missedHeartbeats := range r.heartbeatLostRegistry {

		if missedHeartbeats > 3 {
			// If the connection has missed more than 3 heartbeat_losts, remove it
			logger.Infoln("[ConnectionRegistry]remove connection for heartbeat timeout[🕔💔] address:", client.RemoteAddr())
			r.RemoveClient(client, true)
		} else {
			i := r.heartbeatLostRegistry[client]
			r.heartbeatLostRegistry[client] = i + 1
		}

	}
}

func (r *ClientConnectionRegistry) evictConnectorLost() {
	for client, localClient := range r.readyLocalClients {
		isConnectorStopped := localClient.Connector != nil && (*localClient.Connector).IsAlive()
		if localClient.Stop || isConnectorStopped {
			logger.Infof("[ConnectionRegistry]remove connection for Connector Lost timeout[🕔⏹] address:%s isConnectorStopped: %t", client.RemoteAddr(), isConnectorStopped)
			r.RemoveClient(client, true)
		}
	}

}
