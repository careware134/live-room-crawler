package platform

import (
	"github.com/gorilla/websocket"
	"live-room-crawler/common"
	"live-room-crawler/constant"
	"live-room-crawler/platform/douyin"
)

type IPlatformConnectorStrategy interface {
	GetRoomInfo() *common.RoomInfo

	Connect(localConn *websocket.Conn) constant.ResponseStatus

	StartListen(localConn *websocket.Conn)

	Stop()

	IsAlive() bool
}

func NewConnector(targetStruct common.TargetStruct, stopChan chan struct{}) IPlatformConnectorStrategy {
	if targetStruct.Platform == common.DOUYIN {
		return douyin.NewInstance(targetStruct.LiveURL, stopChan)
	}
	return nil
}
