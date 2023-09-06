package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cors "github.com/rs/cors/wrapper/gin"
	"live-room-crawler/registry/connection"
	"live-room-crawler/util"
	"net/http"
	"strconv"
)

var logger = util.Logger()

func StartLocalServer(port int) {
	// get registry instance
	clientRegistry := connection.GetClientRegistry()
	// dataRegistry := data.GetDataRegistry()

	go clientRegistry.StartHeartbeatsCheck()
	// go dataRegistry.StartPushPlayMessage() // move this block to dataRegistry.MarkReady

	// Define a handler function for WebSocket connections
	// Register the handler function for requests to the "/ws" path
	server := gin.Default()
	server.Use(gin.Recovery())
	server.Use(cors.Default())

	// task handle
	server.GET("/v1", func(c *gin.Context) {
		listenHandler(c.Writer, c.Request)
	})
	addr := ":" + strconv.Itoa(port)
	// Start the server
	logger.Info("[server]Starting server on port:", port)
	if err := server.Run(addr); err != nil {
		logger.Fatal("[server]ListenAndServe Fail with err: ", err)
		return
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
	client := connection.NewClient(conn)
	clientRegistry := connection.GetClientRegistry()
	clientRegistry.AddTempClient(conn)

	go client.StartListenConnection(conn)
}
