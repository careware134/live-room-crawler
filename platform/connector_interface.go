package platform

import (
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/platform/douyin"
	"live-room-crawler/platform/kuaishou"
)

type IPlatformConnectorStrategy interface {
	GetRoomInfo() *domain.RoomInfo

	Connect() constant.ResponseStatus

	StartListen(localConn *websocket.Conn)

	Stop()

	IsAlive() bool
}

func NewConnector(targetStruct domain.TargetStruct, stopChan chan struct{}, localConn *websocket.Conn) IPlatformConnectorStrategy {
	if targetStruct.Platform == domain.DOUYIN {
		return douyin.NewInstance(targetStruct, stopChan, localConn)
	}
	if targetStruct.Platform == domain.KUAISHOU {
		return kuaishou.NewInstance(targetStruct, stopChan, localConn)
	}
	return nil
}
