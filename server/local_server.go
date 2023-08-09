package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"live-room-crawler/registry/connection"
	"live-room-crawler/registry/data"
	"live-room-crawler/util"
	"net/http"
	"strconv"
)

var logger = util.Logger()

func StartLocalServer(port int) {
	// get registry instance
	clientRegistry := connection.GetClientRegistry()
	dataRegistry := data.GetDataRegistry()

	// Define a handler function for WebSocket connections
	// Register the handler function for requests to the "/ws" path
	http.HandleFunc("/v1", listenHandler)

	go clientRegistry.StartHeartbeatsCheck()
	go dataRegistry.StartPushPlayMessage()

	// Start the server
	logger.Info("[server]Starting server on port:", port)
	addr := ":" + strconv.Itoa(port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Fatal("[server]ListenAndServe Fail with err: ", err)
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
	// Add the connection to the registry
	clientRegistry := connection.GetClientRegistry()
	clientRegistry.AddClient(conn)

	// close conn finally
	defer conn.Close()

	client := connection.NewClient(conn)

	// Read messages from the connection
	for {
		select {
		case <-client.StopChan:
			// Stop signal received, exit the goroutine
			return
		default:
			privateStartListen(&client, clientRegistry)
		}
	}
}

func privateStartListen(client *connection.LocalClient, clientRegistry *connection.ClientConnectionRegistry) {
	conn := client.Conn
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		logger.Warn("[server]StartListenConnection FAIL to read message [🆖] error:", err)
		clientRegistry.RemoveClient(conn, true)
		client.TryRevoke()
	}

	// Handle the message
	if messageType == websocket.TextMessage {
		// response to the connection
		response := client.OnCommand(message)
		if response.CommandType == common.STOP && response.ResponseStatus.Success {
			logger.Warnf("[server]🪝StartListenConnection Break by stop request: %s", message)
			client.TryRevoke()
		}

		marshal, err := json.Marshal(response)
		if err != nil {
			logger.Errorf("listenHandler fail to Marshal json response:%e", err)
			client.TryRevoke()
		}
		err = conn.WriteMessage(messageType, marshal)
		if err != nil {
			logger.Errorf("listenHandler fail to WriteMessage for conn:%s error:%e", conn.RemoteAddr(), err)
			client.TryRevoke()
		}
	} else if messageType == websocket.PingMessage {
		// update the heartbeat lost count for this connection
		clientRegistry.UpdateHeartBeat(conn)
		// response with pong
		err = conn.WriteMessage(websocket.PongMessage, []byte("{\"type\":\"pong\"}"))
		if err != nil {
			logger.Errorf("listenHandler fail to WriteMessage Pong for conn:%s error:%e", conn.RemoteAddr(), err)
			client.TryRevoke()
		}
	}
}
