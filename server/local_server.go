package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cors "github.com/rs/cors/wrapper/gin"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/platform"
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

	// task handle
	server.POST("/verifyRoomInfo", func(c *gin.Context) {
		verifyHandler(c.Writer, c.Request)
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

func verifyHandler(w gin.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var target domain.TargetStruct
	var verifyResponse *domain.CommandResponse
	requestJson := r.Body
	logger.Infof("verifyHandler request with: %s", requestJson)
	err := json.NewDecoder(requestJson).Decode(&target)
	if err != nil {
		verifyResponse.ResponseStatus = constant.INVALID_PARAM
		writeVerifyResponse(verifyResponse, w)
		return
	}

	connector := platform.NewConnector(target, nil, nil)
	if connector == nil {
		verifyResponse.ResponseStatus = constant.UNKNOWN_PLATFORM
		writeVerifyResponse(verifyResponse, w)
	}

	verifyResponse = connector.VerifyTarget()
	writeVerifyResponse(verifyResponse, w)
}

func writeVerifyResponse(verifyResponse *domain.CommandResponse, w gin.ResponseWriter) {
	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	response, _ := json.Marshal(verifyResponse)
	w.Write(response)
	logger.Infof("writeVerifyResponse with: %s", response)
}
