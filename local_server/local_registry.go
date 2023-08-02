package local_server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

var reg LocalClientRegistry

type LocalClientRegistry struct {
	m                     sync.Mutex
	clients               map[*websocket.Conn]int
	heartbeatLostRegistry map[*websocket.Conn]int
	readyCommandRegistry  map[*websocket.Conn]CommandRequest
	roomInfo              *RoomInfo
}

func (r *LocalClientRegistry) OnCommand(message []byte) *CommandResponse {
	logger.Infof("[🛎📩]OnCommand with request: %s", string(message))
	request := &CommandRequest{}
	json.Unmarshal(message, request)

	response := &CommandResponse{
		Success: false,
		Message: "unknown",
	}
	switch request.CommandType {
	case START:
		response = tryStart(request)
	case LOAD:
		response = tryLoad(request)
	case STOP:
		response = tryStop(request)
	}

	marshal, _ := json.Marshal(response)
	logger.Infof("[🛎📤]OnCommand return response: %s", marshal)
	return response
}

func tryStart(request *CommandRequest) *CommandResponse {
	response := &CommandResponse{
		CommandType: START,
		TraceId:     request.TraceId,
		Success:     true,
		Message:     "start watch success",
	}
	return response
}

func tryLoad(request *CommandRequest) *CommandResponse {
	response := &CommandResponse{
		CommandType: LOAD,
		TraceId:     request.TraceId,
		Success:     true,
		Message:     "load project rule success",
	}
	return response
}

func tryStop(request *CommandRequest) *CommandResponse {
	response := &CommandResponse{
		CommandType: STOP,
		TraceId:     request.TraceId,
		Success:     true,
		Message:     "stop watch success",
	}
	return response
}

//func TryPong(request *CommandRequest) *CommandResponse {
//	response := &CommandResponse{
//		CommandType: PONG,
//		TraceId:     request.TraceId,
//		Success:     true,
//		Message:     "success",
//	}
//	return response
//}
