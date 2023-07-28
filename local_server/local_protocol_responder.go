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

func (r *LocalClientRegistry) ResponseClient(message []byte) *CommandResponse {
	request := &CommandRequest{}
	json.Unmarshal(message, request)
	if request.CommandType == START {

	}

	return &CommandResponse{}
}
