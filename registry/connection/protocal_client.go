package connection

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/platform"
	"live-room-crawler/registry/data"
	"sync"

	//"live-room-crawler/registry/connection"
	"live-room-crawler/util"
)

var logger1 = util.Logger()

type LocalClient struct {
	Connector *platform.IPlatformConnectorStrategy
	Conn      *websocket.Conn
	Start     bool
	Stop      bool
	stopOnce  sync.Once
	stopChan  chan struct{} // Channel to signal stop
}

func NewClient(conn *websocket.Conn) LocalClient {
	client := LocalClient{
		Conn:     conn,
		stopChan: make(chan struct{}),
	}
	client.setPingHandler()
	return client
}

func (client *LocalClient) StartListenConnection(conn *websocket.Conn) {
	logger.Infof("[LocalClient]▶️StartListenConnection for client: %s", conn.RemoteAddr())

	// close conn finally
	defer conn.Close()

	clientRegistry := GetClientRegistry()
	dataRegistry := data.GetDataRegistry()

	// Read messages from the connection
	for {
		select {
		case <-client.stopChan:
			// Stop signal received, exit the goroutine
			logger.Infof("⚓️[LocalClient]break StartListenConnection by client.stopChan for client: %s", conn.RemoteAddr())
			return
		default:
			client.privateStartListen(clientRegistry, dataRegistry)
		}
	}
}

func (client *LocalClient) privateStartListen(clientRegistry *ClientConnectionRegistry, dataRegistry *data.EventDataRegistry) {
	conn := client.Conn
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		logger.Warn("[LocalClient]StartListenConnection FAIL to read message [🆖] error:", err)
		clientRegistry.RemoveClient(conn, true)
		return
	}

	// Handle the message
	if messageType == websocket.TextMessage {
		// response to the connection
		response := client.OnCommand(message)
		if response.CommandType == domain.STOP && response.ResponseStatus.Success {
			logger.Warnf("[LocalClient]🪝StartListenConnection Break by stop request: %s", message)
			client.TryRevoke()
			return
		}

		err = dataRegistry.WriteResponse(conn, response)
		if err != nil {
			logger.Errorf("[LocalClient]listenHandler fail to WriteResponse for conn:%s error:%e", conn.RemoteAddr(), err)
			client.TryRevoke()
			return
		}
	}
}

func (client *LocalClient) OnCommand(
	message []byte) *domain.CommandResponse {

	logger1.Infof("[🛎📩LocalClient]OnCommand request is: %s", string(message))
	request := &domain.CommandRequest{}
	json.Unmarshal(message, request)

	response := &domain.CommandResponse{
		CommandType:    request.CommandType,
		ResponseStatus: constant.UNKNOWN_COMMAND,
	}
	switch request.CommandType {
	case domain.START:
		response = client.onStart(request)
	//case domain.LOAD:
	//	response = client.onRefresh(request)
	case domain.REFRESH:
		response = client.onRefresh(request)
	case domain.STOP:
		response = client.onStop(request)
	case domain.PING:
		GetClientRegistry().UpdateHeartBeat(client.Conn)
		response = &domain.CommandResponse{
			CommandType: domain.PONG,
		}
	}

	marshal, _ := json.Marshal(response)
	logger1.Infof("[🛎📤LocalClient]OnCommand response is: %s", marshal)
	return response
}

func (client *LocalClient) setPingHandler() {
	clientRegistry := GetClientRegistry()
	dataRegistry := data.GetDataRegistry()
	conn := client.Conn
	conn.SetPingHandler(func(appData string) error {
		clientRegistry.UpdateHeartBeat(conn)
		// response with pong
		err := dataRegistry.WriteResponse(conn, &domain.CommandResponse{
			CommandType:    domain.PONG,
			ResponseStatus: constant.SUCCESS,
		})
		if err != nil {
			logger.Errorf("PingHandler fail to WriteMessage Pong for conn:%s error:%e", conn.RemoteAddr(), err)
			//client.TryRevoke()
		}
		return nil
	})
}

func (client *LocalClient) onStart(request *domain.CommandRequest) *domain.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger1.Infof("🌏onStart with request: %s", marshal)

	response := &domain.CommandResponse{
		CommandType:    domain.START,
		TraceId:        request.TraceId,
		ResponseStatus: constant.SUCCESS,
	}

	// ensure idempotent request
	if data.GetDataRegistry().IsReady(client.Conn) {
		response.ResponseStatus = constant.SUCCESS_ALREADY
		return response
	}

	// create connector by start request
	connector := platform.NewConnector(request.Target, client.stopChan)
	if connector == nil {
		response.ResponseStatus = constant.UNKNOWN_PLATFORM
		logger.Errorf("onStart fail with UNKNOWN_PLATFORM")
		return response
	}
	client.Connector = &connector
	//info := connector.GetRoomInfo()
	//if info == nil {
	//	response.ResponseStatus = constant.INVALID_LIVE_URL
	//	logger.Errorf("onStart fail with INVALID_LIVE_URL")
	//	return response
	//}

	// 0. invoke connect to prepare listen
	responseStatus := connector.Connect()
	if !responseStatus.Success {
		response.ResponseStatus = responseStatus
		return response
	}
	response.Room = *connector.GetRoomInfo()

	// .1 then start mark ready
	GetClientRegistry().MarkReady(client.Conn, request, client)

	// .2 if connect success, then load rule
	responseStatus = data.GetDataRegistry().LoadRule(request.TraceId, client.Conn)
	if !responseStatus.Success {
		response.ResponseStatus = responseStatus
		GetClientRegistry().RemoveClient(client.Conn, false)
		return response
	}
	// .3 start listen connector
	go connector.StartListen(client.Conn)
	// .4 mark connection start
	client.Start = true
	client.Stop = false

	marshal, _ = json.Marshal(response)
	logger1.Infof("🌏onStart with response: %s", marshal)

	return response
}

func (client *LocalClient) onRefresh(request *domain.CommandRequest) *domain.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger1.Infof("🌏onRefresh with request: %s", marshal)

	responseStatus := constant.SUCCESS

	response := &domain.CommandResponse{
		CommandType:    domain.REFRESH,
		TraceId:        request.TraceId,
		ResponseStatus: responseStatus,
	}

	// .2 if connect success, then load rule
	responseStatus = data.GetDataRegistry().LoadRule(request.TraceId, client.Conn)
	if !responseStatus.Success {
		response.ResponseStatus = responseStatus
		return response
	}

	marshal, _ = json.Marshal(response)
	logger1.Infof("🌏onRefresh with response: %s", marshal)

	return response
}

func (client *LocalClient) onStop(request *domain.CommandRequest) *domain.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger1.Infof("🌏onStop with request: %s", marshal)

	response := &domain.CommandResponse{
		CommandType:    domain.STOP,
		TraceId:        request.TraceId,
		ResponseStatus: constant.SUCCESS,
	}

	client.stopOnce.Do(func() {
		client.privateTryStop(response)
		GetClientRegistry().RemoveClient(client.Conn, false)
	})

	marshal, _ = json.Marshal(response)
	logger1.Infof("🌏onStop with response: %s", marshal)
	return response
}

func (client *LocalClient) TryRevoke() *domain.CommandResponse {
	response := &domain.CommandResponse{
		CommandType:    domain.STOP,
		TraceId:        "revoke-" + uuid.NewString(),
		ResponseStatus: constant.LIVE_CONNECTION_CLOSED,
	}

	marshal, _ := json.Marshal(response)
	logger1.Infof("🌏TryRevoke with response: %s", marshal)

	client.stopOnce.Do(func() {
		client.privateTryStop(response)
	})

	return response
}

func (client *LocalClient) privateTryStop(response *domain.CommandResponse) {
	data.GetDataRegistry().WriteResponse(client.Conn, response)
	client.Conn.Close()

	if !client.Stop {
		close(client.stopChan)
		client.Stop = true
		client.Start = false
	}

	if client.Connector != nil {
		(*client.Connector).Stop()
	}
	logger.Infof("privateTryStop with client:%s", client.Conn.RemoteAddr())
}
