package platform

import (
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/platform/douyin"
	"live-room-crawler/platform/kuaishou"
	"live-room-crawler/platform/wheadless"
)

type IPlatformConnectorStrategy interface {
	GetRoomInfo() *domain.RoomInfo

	SetRoomInfo(info domain.RoomInfo)

	Connect() constant.ResponseStatus

	StartListen(localConn *websocket.Conn)

	Stop()

	IsDead() bool

	VerifyTarget() *domain.CommandResponse
}

func NewConnector(targetStruct domain.TargetStruct, stopChan chan struct{}, localConn *websocket.Conn) IPlatformConnectorStrategy {

	if targetStruct.Headless {
		return wheadless.NewInstance(targetStruct, stopChan, localConn)
	}

	if targetStruct.Platform == domain.DOUYIN {
		return douyin.NewInstance(targetStruct, stopChan, localConn)
	}
	if targetStruct.Platform == domain.KUAISHOU {
		return kuaishou.NewInstance(targetStruct, stopChan, localConn)
	}
	return nil
}
