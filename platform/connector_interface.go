package platform

import (
	"live-room-crawler/common"
	"live-room-crawler/constant"
	"live-room-crawler/platform/douyin"
)

type IPlatformConnectorStrategy interface {
	GetRoomInfo() *common.RoomInfo

	Start() constant.ResponseStatus

	Stop()
}

func NewConnector(targetStruct common.TargetStruct) IPlatformConnectorStrategy {
	if targetStruct.Platform == common.DOUYIN {
		return douyin.NewInstance(targetStruct.LiveURL)
	}
	return nil
}
