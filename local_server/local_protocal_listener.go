package local_server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"sync"
)

type LocalClientRegistry struct {
	m                     sync.Mutex
	clients               map[*websocket.Conn]int
	heartbeatLostRegistry map[*websocket.Conn]int
	readyCommandRegistry  map[*websocket.Conn]common.CommandRequest
	roomInfo              *common.RoomInfo
}

func (r *LocalClientRegistry) OnCommand(conn *websocket.Conn, message []byte) *common.CommandResponse {
	logger.Infof("[🛎📩]OnCommand request is: %s", string(message))
	request := &common.CommandRequest{}
	json.Unmarshal(message, request)

	response := &common.CommandResponse{
		CommandType: request.CommandType,
		ResponseStatus: common.ResponseStatus{
			Success: false,
			Message: "unknown",
		},
	}
	switch request.CommandType {
	case common.START:
		response = r.tryStart(conn, request)
	case common.LOAD:
		response = r.tryLoad(conn, request)
	case common.STOP:
		response = r.tryStop(conn, request)
	}

	marshal, _ := json.Marshal(response)
	logger.Infof("[🛎📤]OnCommand response is: %s", marshal)
	return response
}

func (r *LocalClientRegistry) tryStart(conn *websocket.Conn, request *common.CommandRequest) *common.CommandResponse {
	response := &common.CommandResponse{
		CommandType: common.START,
		TraceId:     request.TraceId,
		ResponseStatus: common.ResponseStatus{
			Success: true,
			Message: "start watch success",
		},
	}
	r.markReady(conn, *request)
	return response
}

func (r *LocalClientRegistry) tryLoad(conn *websocket.Conn, request *common.CommandRequest) *common.CommandResponse {
	response := &common.CommandResponse{
		CommandType: common.LOAD,
		TraceId:     request.TraceId,
		ResponseStatus: common.ResponseStatus{
			Success: true,
			Message: "load project rule success",
		},
	}
	return response
}

func (r *LocalClientRegistry) tryStop(conn *websocket.Conn, request *common.CommandRequest) *common.CommandResponse {
	response := &common.CommandResponse{
		CommandType: common.STOP,
		TraceId:     request.TraceId,
		ResponseStatus: common.ResponseStatus{
			Success: true,
			Message: "stop watch success",
		},
	}
	r.removeClient(conn)
	return response
}

func (r *LocalClientRegistry) Broadcast(updateRegistryEvent *common.UpdateRegistryEvent) {
	r.m.Lock()
	defer r.m.Unlock()
	for client := range r.readyCommandRegistry {
		response := &common.CommandResponse{
			CommandType: common.PLAY,
			TraceId:     uuid.NewString(),
			RuleMeta: common.RuleMeta{
				Id:   1,
				Name: "mock-测试规则",
			},
			Content: common.PlayContent{
				DrivenType: common.TEXT,
				Text:       updateRegistryEvent.ActionList[0].Content,
			},
		}
		marshal, err := json.Marshal(response)
		if err != nil {
			logger.Info("Broadcast error Marshal:", err)
			continue
		}
		err = client.WriteMessage(websocket.TextMessage, marshal)
		if err != nil {
			logger.Info("Broadcast error WriteMessage:", err)
			r.removeClient(client)
		}
	}
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
