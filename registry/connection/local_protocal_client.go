package connection

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"live-room-crawler/constant"
	"live-room-crawler/platform"
	//"live-room-crawler/registry/connection"
	"live-room-crawler/util"
)

var logger1 = util.Logger()

type LocalClient struct {
	Connector *platform.IPlatformConnectorStrategy
	Conn      *websocket.Conn
	Start     bool
	Stop      bool
}

func NewClient(conn *websocket.Conn) LocalClient {
	return LocalClient{
		Conn: conn,
	}
}

func (client *LocalClient) OnCommand(
	message []byte) *common.CommandResponse {

	logger1.Infof("[🛎📩]OnCommand request is: %s", string(message))
	request := &common.CommandRequest{}
	json.Unmarshal(message, request)

	response := &common.CommandResponse{
		CommandType: request.CommandType,
	}
	switch request.CommandType {
	case common.START:
		response = client.tryStart(request)
	case common.LOAD:
		response = client.tryLoad(request)
	case common.STOP:
		response = client.TryStop(request)
	}

	marshal, _ := json.Marshal(response)
	logger1.Infof("[🛎📤]OnCommand response is: %s", marshal)
	return response
}

func (client *LocalClient) tryStart(request *common.CommandRequest) *common.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger1.Infof("🌏tryStart with request: %s", marshal)

	response := &common.CommandResponse{
		CommandType: common.START,
		TraceId:     request.TraceId,
	}

	// create connector by start request
	connector := platform.NewConnector(request.Target)
	client.Connector = &connector
	// invoke connect
	responseStatus := connector.Connect(client.Conn)
	response.ResponseStatus = responseStatus
	if !responseStatus.Success {
		return response
	}

	// if connect success,
	// .1 then start mark ready
	GetClientRegistry().MarkReady(client.Conn, request, client)
	// .2 start listen
	go connector.StartListen(client.Conn)
	// .3 mark connection start
	client.Start = true

	marshal, _ = json.Marshal(response)
	logger1.Infof("🌏tryStart with response: %s", marshal)

	return response
}

func (client *LocalClient) tryLoad(request *common.CommandRequest) *common.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger1.Infof("🌏tryLoad with request: %s", marshal)

	response := &common.CommandResponse{
		CommandType: common.LOAD,
		TraceId:     request.TraceId,
	}

	response.ResponseStatus = constant.SUCCESS
	marshal, _ = json.Marshal(response)
	logger1.Infof("🌏tryLoad with response: %s", marshal)

	return response
}

func (client *LocalClient) TryStop(request *common.CommandRequest) *common.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger1.Infof("🌏TryStop with request: %s", marshal)

	TraceId := request.TraceId
	Message := constant.SUCCESS.Message

	response := &common.CommandResponse{
		CommandType: common.STOP,
		TraceId:     TraceId,
		ResponseStatus: constant.ResponseStatus{
			Success: true,
			Message: Message,
		},
	}

	client.privateTryStop(response)
	GetClientRegistry().RemoveClient(client.Conn, false)

	marshal, _ = json.Marshal(response)
	logger1.Infof("🌏TryStop with response: %s", marshal)
	return response
}

func (client *LocalClient) TryRevoke(request *common.CommandRequest) *common.CommandResponse {
	response := &common.CommandResponse{
		CommandType:    common.STOP,
		TraceId:        "revoke-" + uuid.NewString(),
		ResponseStatus: constant.LIVE_CONNECTION_CLOSED,
	}

	marshal, _ := json.Marshal(response)
	logger1.Infof("🌏TryRevoke with response: %s", marshal)

	client.privateTryStop(response)
	return response
}

func (client *LocalClient) privateTryStop(response *common.CommandResponse) {
	marshal, _ := json.Marshal(response)
	client.Conn.WriteMessage(websocket.TextMessage, marshal)
	client.Start = false
	client.Stop = true
	client.Conn.Close()

	(*client.Connector).Stop()
}
