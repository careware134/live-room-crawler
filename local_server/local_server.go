package local_server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"live-room-crawler/registry"
	"live-room-crawler/util"
	"net/http"
	"strconv"
)

var logger = util.Logger()

func StartLocalServer(port int) {
	// get registry instance
	clientRegistry := registry.GetInstance()

	// Define a handler function for WebSocket connections
	// Register the handler function for requests to the "/ws" path
	http.HandleFunc("/v1", listenHandler)

	go clientRegistry.StartHeartbeatsCheck()
	go clientRegistry.StartPushPlayMessage()

	// Start the local_server
	logger.Info("[local_server]Starting local_server on port:", port)
	addr := ":" + strconv.Itoa(port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Fatal("[local_server]ListenAndServe Fail with err: ", err)
	}
}

func listenHandler(response http.ResponseWriter, request *http.Request) {
	// Define the WebSocket upgrade function
	upgrade := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	conn, err := upgrade.Upgrade(response, request, nil)
	if err != nil {
		logger.Println(err)
		return
	}

	go StartListenConnection(conn)
}

func StartListenConnection(conn *websocket.Conn) {
	// Add the client to the registry
	clientRegistry := registry.GetInstance()
	clientRegistry.AddClient(conn)

	// close conn finally
	defer conn.Close()

	client := NewClient(*conn)

	// Read messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logger.Warn("[local_server]StartListenConnection FAIL to read message [🆖] error:", err)
			clientRegistry.RemoveClient(conn)
			break
		}

		clientRegistry.HeartbeatLostRegistry[conn] = 0

		// Handle the message
		if messageType == websocket.TextMessage {
			clientRegistry.HeartbeatLostRegistry[conn] = 0
			// response to the client
			response := client.OnCommand(conn, message)
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
			clientRegistry.HeartbeatLostRegistry[conn] = 0
			// response with pong
			err = conn.WriteMessage(websocket.PongMessage, []byte("{\"type\":\"pong\"}"))
			if err != nil {
				logger.Errorf("listenHandler fail to WriteMessage Pong for conn:%s error:%e", conn.RemoteAddr(), err)
				continue
			}
		}
	}
}
