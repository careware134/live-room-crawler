package local_server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"live-room-crawler/util"
	"net/http"
	"strconv"
	"time"
)

var logger = util.Logger()

func StartLocalServer(port int) LocalClientRegistry {
	// Create a new registry
	reg := LocalClientRegistry{
		clients:               make(map[*websocket.Conn]int),
		heartbeatLostRegistry: make(map[*websocket.Conn]int),
		readyCommandRegistry:  make(map[*websocket.Conn]common.CommandRequest),
	}

	// Define a handler function for WebSocket connections
	// Register the handler function for requests to the "/ws" path
	http.HandleFunc("/v1", reg.listenHandler)

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

func (r *LocalClientRegistry) markReady(client *websocket.Conn, startRequest common.CommandRequest) {
	r.m.Lock()
	defer r.m.Unlock()
	r.readyCommandRegistry[client] = startRequest
}

func (r *LocalClientRegistry) removeClient(client *websocket.Conn) {
	r.m.Lock()
	defer r.m.Unlock()
	logger.Info("removeClient invoke client addr:", client.RemoteAddr())
	delete(r.clients, client)
	delete(r.heartbeatLostRegistry, client)
	delete(r.readyCommandRegistry, client)
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

func (r *LocalClientRegistry) listenHandler(response http.ResponseWriter, request *http.Request) {
	// Define the WebSocket upgrade function
	upgrade := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	conn, err := upgrade.Upgrade(response, request, nil)
	if err != nil {
		logger.Println(err)
		return
	}

	go r.handleConnection(conn)
}

func (r *LocalClientRegistry) handleConnection(conn *websocket.Conn) {
	// Add the client to the registry
	r.addClient(conn)
	// close conn finally
	defer conn.Close()

	// Read messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Warn("[local_server]read message [🆖] error:", err)
			r.removeClient(conn)
			continue
		}

		r.heartbeatLostRegistry[conn] = 0

		// Handle the message
		if messageType == websocket.TextMessage {
			r.heartbeatLostRegistry[conn] = 0
			// response to the client
			response := r.OnCommand(conn, message)
			marshal, err := json.Marshal(response)
			if err != nil {
				logger.Error("listenHandler fail to Marshal json response:", err)
				continue
			}
			err = conn.WriteMessage(messageType, marshal)
			if err != nil {
				logger.Errorf("listenHandler fail to WriteMessage for conn:%s error:%e", conn.RemoteAddr(), err)
				continue
			}
		} else if messageType == websocket.PingMessage {
			// update the heartbeat lost count for this client
			r.heartbeatLostRegistry[conn] = 0
			// response with pong
			err = conn.WriteMessage(websocket.PongMessage, []byte("{\"type\":\"pong\"}"))
			if err != nil {
				logger.Errorf("listenHandler fail to WriteMessage Pong for conn:%s error:%e", conn.RemoteAddr(), err)
				continue
			}
		}
	}
}
