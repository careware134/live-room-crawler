package local_server

import (
	"github.com/gorilla/websocket"
	"live-room-crawler/util"
	"net/http"
	"strconv"
	"time"
)

var logger = util.Logger()

func StartLocalServer(port int, roomInfo *RoomInfo) LocalClientRegistry {
	// Create a new registry
	reg = LocalClientRegistry{
		clients:               make(map[*websocket.Conn]int),
		heartbeatLostRegistry: make(map[*websocket.Conn]int),
		readyCommandRegistry:  make(map[*websocket.Conn]CommandRequest),
		roomInfo:              roomInfo,
	}

	// Define a handler function for WebSocket connections
	// Register the handler function for requests to the "/ws" path
	http.HandleFunc("/", listenHandler)

	go reg.checkHeartbeats()

	// Start the local_server
	logger.Info("[local_server]Starting local_server on port:", port)
	addr := ":" + strconv.Itoa(port)
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			logger.Fatal("[local_server]ListenAndServe Fail with err: ", err)
		}
	}()

	return reg
}

func (r *LocalClientRegistry) addClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	r.clients[client] = 0
	r.heartbeatLostRegistry[client] = 0
}

func (r *LocalClientRegistry) removeClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	logger.Info("removeClient invoke client addr:", client.RemoteAddr())
	delete(r.clients, client)
	delete(r.heartbeatLostRegistry, client)
}

func (r *LocalClientRegistry) Broadcast(message []byte) {
	r.m.Lock()
	defer r.m.Unlock()
	for client := range r.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logger.Println(err)
			r.removeClient(client)
		}
	}
}

func (r *LocalClientRegistry) checkHeartbeats() {
	for {
		logger.Info("[local_server]start checkHeartbeats[❤️]....")
		// Sleep for 10 seconds
		time.Sleep(10 * time.Second)

		// Check each connected client
		for client, missedHeartbeats := range r.heartbeatLostRegistry {
			if missedHeartbeats > 3 {
				// If the client has missed more than 3 heartbeat_losts, remove it
				logger.Infoln("[local_server]remove client for heartbeat timeout[🕔💔] address:", client.RemoteAddr())
				r.removeClient(client)
			} else {
				// Otherwise, send a heartbeat and increment the missed heartbeat count
				err := client.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					logger.Error("[local_server]fail to write pong to client", err)
					r.removeClient(client)
				} else {
					r.heartbeatLostRegistry[client]++
				}
			}
		}
	}
}

func listenHandler(w http.ResponseWriter, r *http.Request) {
	// Define the WebSocket upgrade function
	upgrade := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		logger.Println(err)
		return
	}

	// Add the client to the registry
	reg.addClient(conn)

	// Read messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Warn("[local_server]read message [🆖] error:", err)
			reg.removeClient(conn)
			continue
		}

		reg.heartbeatLostRegistry[conn] = 0

		// Handle the message
		if messageType == websocket.TextMessage {
			reg.heartbeatLostRegistry[conn] = 0
			// response to the client
			reg.ResponseClient(message)
		} else if messageType == websocket.PingMessage {
			// update the heartbeat losts count for this client
			reg.heartbeatLostRegistry[conn] = 0
		}
	}

}
