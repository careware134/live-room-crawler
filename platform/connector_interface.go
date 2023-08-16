package platform

import (
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/platform/douyin"
)

type IPlatformConnectorStrategy interface {
	GetRoomInfo() *domain.RoomInfo

	Connect(localConn *websocket.Conn) constant.ResponseStatus

	StartListen(localConn *websocket.Conn)

	Stop()

	IsAlive() bool
}

func NewConnector(targetStruct domain.TargetStruct, stopChan chan struct{}) IPlatformConnectorStrategy {
	if targetStruct.Platform == domain.DOUYIN {
		return douyin.NewInstance(targetStruct.LiveURL, stopChan)
	}
	return nil
}
