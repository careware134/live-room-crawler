package kuaishou

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/platform/kuaishou/kuaishou_protostub"
	"live-room-crawler/registry/data"
	"live-room-crawler/util"
	"net/http"
)

var (
	log = util.Logger()
)

type ConnectorStrategy struct {
	Target        domain.TargetStruct
	RoomInfo      *domain.RoomInfo
	conn          *websocket.Conn
	localConn     *websocket.Conn
	IsStart       bool
	IsStop        bool
	stopChan      chan struct{}
	ExtensionInfo ExtensionInfo
}

type ExtensionInfo struct {
	liveUrl string
	//token        string
	cookie string
	//webSocketUrl string
	liveRoomId string
	Headers    http.Header
}

var (
	logger = util.Logger()
)

func NewInstance(Target domain.TargetStruct, stopChan chan struct{}, localConn *websocket.Conn) *ConnectorStrategy {
	util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] NewInstance for url: %s cookie:%s", Target.LiveURL, Target.Cookie)
	return &ConnectorStrategy{
		Target:    Target,
		localConn: localConn,
		stopChan:  stopChan,
	}
}

func (c *ConnectorStrategy) SetRoomInfo(info domain.RoomInfo) {
	marshal, _ := json.Marshal(info)
	c.RoomInfo = &info
	util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] SetRoomInfo with value:%s", marshal)
}

func (c *ConnectorStrategy) VerifyTarget() *domain.CommandResponse {
	info := c.GetRoomInfo()
	responseStatus := constant.SUCCESS
	if info == nil {
		responseStatus = constant.FAIL_GETTING_ROOM_INFO
		return &domain.CommandResponse{
			ResponseStatus: responseStatus,
		}
	}

	_, err := c.GetWebSocketInfo(info.RoomId)
	if err != nil {
		responseStatus = constant.FAIL_GETTING_SOCKET_INFO
	}

	return &domain.CommandResponse{
		Room:           *info,
		ResponseStatus: responseStatus,
	}
}

func (connector *ConnectorStrategy) Connect() constant.ResponseStatus {
	roomInfo := connector.GetRoomInfo()
	if roomInfo == nil {
		util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] fail to get kuaishou roomInfo with url: %s", connector.Target.LiveURL)
		return constant.FAIL_GETTING_ROOM_INFO
	}
	if roomInfo.WebSocketUrl == "" {
		return constant.FAIL_GETTING_SOCKET_INFO
	}

	header := connector.ExtensionInfo.Headers
	url := connector.RoomInfo.WebSocketUrl
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		util.Logger().Infof("[kuaishou.connnector][connect]fail to dial to addr:%s ", connector.Target.LiveURL)
		return constant.CONNECTION_FAIL
	}
	connector.conn = conn
	util.Logger().Infof("[kuaishou.connnector][connect]dial to addr:%s ", connector.Target.LiveURL)

	obj := kuaishou_protostub.CSWebEnterRoom{
		PayloadType: 200,
		Payload: &kuaishou_protostub.CSWebEnterRoom_Payload{
			Token:        connector.RoomInfo.Token,
			LiveStreamId: connector.RoomInfo.RoomId,
			PageId:       connector.getPageID(),
		},
	}

	marshal, err := proto.Marshal(&obj)
	if err != nil {
		util.Logger().Errorf("[kuaishou.connector][ConnectData] [序列化失败] with error: %e", err)
		return constant.CONNECTION_FAIL
	}

	util.Logger().Infof("[kuaishou.connector][OnOpen] [建立wss连接] with client: %s", conn.RemoteAddr())
	err = conn.WriteMessage(websocket.BinaryMessage, marshal)
	if err != nil {
		util.Logger().Infof("[kuaishou.connector][OnOpen] [发送数据失败] with error: %e", err)
		return constant.CONNECTION_FAIL
	}

	go connector.StartHeartbeat(conn)
	return constant.SUCCESS
}

func (connector *ConnectorStrategy) GetRoomInfo() *domain.RoomInfo {
	if connector.RoomInfo != nil {
		marshal, _ := json.Marshal(connector.RoomInfo)
		util.Logger().Infof("[kuaishou.connnector]GetRoomInfo SKIP for ALREADY have value:%s", marshal)
		return connector.RoomInfo
	}

	liveRoomId, liveRoomCaption, err := connector.GetLiveRoomId()
	if err != nil {
		util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] fail to get kuaishou roomInfo with url: %s", connector.Target.LiveURL)
		return nil
	}

	connector.RoomInfo = &domain.RoomInfo{
		RoomId:    liveRoomId,
		RoomTitle: liveRoomCaption,
	}

	util.Logger().Infof("[kuaishou.connnector][GetRoomInfo]return with roomId:%s roomTitle:%s", liveRoomId, liveRoomCaption)
	_, err = connector.GetWebSocketInfo(liveRoomId)
	if err != nil {
		util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] fail to get GetWebSocketInfo roomInfo with url: %s", connector.Target.LiveURL)
		return connector.RoomInfo
	}
	return connector.RoomInfo
}

func (connector *ConnectorStrategy) StartListen(localConn *websocket.Conn) {
	util.Logger().Infof("[kuaishou.connnector]StartListen for room:%s", connector.RoomInfo.RoomId)
	connector.IsStart = true
	dataRegistry := data.GetDataRegistry()
	for {
		_, message, err := connector.conn.ReadMessage()

		select {
		case <-connector.stopChan:
			// Stop signal received, exit the goroutine
			util.Logger().Infof("⚓️♪ [kuaishou.ConnectorStrategy] StartListen BREAKED by c.stopChan")
			return
		default:
			if err != nil {
				util.Logger().Errorf("[kuaishou.ConnectorStrategy] StartListen fail with reason: %e", err)
				connector.Stop()
				return
			}
			connector.OnMessage(message, localConn, dataRegistry)
		}
	}

}

func (connector *ConnectorStrategy) Stop() {
	connector.IsStart = false
	connector.IsStop = true
	if connector.conn != nil {
		connector.conn.Close()
	}
	title := ""
	if connector.RoomInfo != nil {
		title = connector.RoomInfo.RoomTitle
	}
	util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] Stop kuaishou for url: %s title: %s", connector.Target.LiveURL, title)
}

func (connector *ConnectorStrategy) IsAlive() bool {
	return connector.IsStop
}
