package local_server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"live-room-crawler/constant"
	"live-room-crawler/platform"
	"live-room-crawler/registry"
)

type LocalClient struct {
	connector platform.IPlatformConnectorStrategy
	conn      websocket.Conn
	Start     bool
}

func NewClient(conn websocket.Conn) LocalClient {
	return LocalClient{
		conn: conn,
	}
}

func (client *LocalClient) OnCommand(
	conn *websocket.Conn,
	message []byte) *common.CommandResponse {

	logger.Infof("[🛎📩]OnCommand request is: %s", string(message))
	request := &common.CommandRequest{}
	json.Unmarshal(message, request)

	response := &common.CommandResponse{
		CommandType: request.CommandType,
	}
	switch request.CommandType {
	case common.START:
		response = client.tryStart(conn, request)
	case common.LOAD:
		response = client.tryLoad(conn, request)
	case common.STOP:
		response = client.tryStop(conn, request)
	}

	marshal, _ := json.Marshal(response)
	logger.Infof("[🛎📤]OnCommand response is: %s", marshal)
	return response
}

func (client *LocalClient) tryStart(conn *websocket.Conn, request *common.CommandRequest) *common.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger.Infof("🌏tryStart with request: %s", marshal)

	response := &common.CommandResponse{
		CommandType: common.START,
		TraceId:     request.TraceId,
	}

	connector := platform.NewConnector(request.Target)
	responseStatus := connector.Start()
	response.ResponseStatus = responseStatus
	if !responseStatus.Success {
		return response
	}

	client.Start = true
	registry.GetInstance().MarkReady(conn, *request, connector)
	marshal, _ = json.Marshal(response)
	logger.Infof("🌏tryStart with response: %s", marshal)

	return response
}

func (client *LocalClient) tryLoad(conn *websocket.Conn, request *common.CommandRequest) *common.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger.Infof("🌏tryLoad with request: %s", marshal)

	response := &common.CommandResponse{
		CommandType: common.LOAD,
		TraceId:     request.TraceId,
	}

	client.Start = false
	response.ResponseStatus = constant.SUCCESS
	marshal, _ = json.Marshal(response)
	logger.Infof("🌏tryLoad with response: %s", marshal)

	return response
}

func (client *LocalClient) tryStop(conn *websocket.Conn, request *common.CommandRequest) *common.CommandResponse {
	marshal, _ := json.Marshal(request)
	logger.Infof("🌏tryStop with request: %s", marshal)

	response := &common.CommandResponse{
		CommandType: common.STOP,
		TraceId:     request.TraceId,
		ResponseStatus: constant.ResponseStatus{
			Success: true,
			Message: "stop watch success",
		},
	}

	marshal, _ = json.Marshal(response)
	logger.Infof("🌏tryStop with response: %s", marshal)
	conn.WriteMessage(websocket.TextMessage, marshal)

	registry.GetInstance().RemoveClient(conn)
	return response
}
