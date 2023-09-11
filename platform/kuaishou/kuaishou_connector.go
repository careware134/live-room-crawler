package kuaishou

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spyzhov/ajson"
	"google.golang.org/protobuf/proto"
	"io"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/platform/kuaishou/kuaishou_protostub"
	"live-room-crawler/registry/data"
	"live-room-crawler/util"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

var (
	log = util.Logger()
)

type ConnectorStrategy struct {
	liveUrl       string
	RoomInfo      *domain.RoomInfo
	conn          *websocket.Conn
	localConn     *websocket.Conn
	IsStart       bool
	IsStop        bool
	stopChan      chan struct{}
	ExtensionInfo ExtensionInfo
}

type ExtensionInfo struct {
	token        string
	cookie       string
	webSocketUrl string
	liveRoomId   string
	liveUrl      string
	Headers      http.Header
}

var (
	logger = util.Logger()
)

func NewInstance(liveUrl string, stopChan chan struct{}) *ConnectorStrategy {
	logger.Infof("🎦[kuaishou.ConnectorStrategy] NewInstance for url: %s", liveUrl)
	return &ConnectorStrategy{
		liveUrl:  liveUrl,
		stopChan: stopChan,
	}
}

func (connector *ConnectorStrategy) Connect() constant.ResponseStatus {
	roomInfo := connector.GetRoomInfo()
	if roomInfo == nil {
		logger.Infof("🎦[kuaishou.ConnectorStrategy] fail to get kuaishou roomInfo with url: %s", connector.liveUrl)
		return constant.FAIL_GETTING_ROOM_INFO
	}

	_, err := connector.GetWebSocketInfo(roomInfo.RoomId)
	if err != nil {
		return constant.CONNECTION_FAIL
	}

	header := connector.ExtensionInfo.Headers
	url := connector.ExtensionInfo.webSocketUrl
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		logger.Infof("[kuaishou.connnector][connect]fail to dial to addr:%s ", url)
		return constant.CONNECTION_FAIL
	}
	connector.conn = conn
	logger.Infof("[kuaishou.connnector][connect]dial to addr:%s ", url)

	obj := kuaishou_protostub.CSWebEnterRoom{
		PayloadType: 200,
		Payload: &kuaishou_protostub.CSWebEnterRoom_Payload{
			Token:        connector.ExtensionInfo.token,
			LiveStreamId: connector.ExtensionInfo.liveRoomId,
			PageId:       connector.getPageID(),
		},
	}

	marshal, err := proto.Marshal(&obj)
	if err != nil {
		logger.Errorf("[kuaishou.connector][ConnectData] [序列化失败] with error: %e", err)
		return constant.CONNECTION_FAIL
	}

	logger.Infof("[kuaishou.connector][OnOpen] [建立wss连接] with client: %s", conn.RemoteAddr())
	err = conn.WriteMessage(websocket.BinaryMessage, marshal)
	if err != nil {
		logger.Infof("[kuaishou.connector][OnOpen] [发送数据失败] with error: %e", err)
		return constant.CONNECTION_FAIL
	}

	go connector.StartHeartbeat(conn)
	return constant.SUCCESS
}

func (connector *ConnectorStrategy) GetRoomInfo() *domain.RoomInfo {
	if connector.RoomInfo != nil {
		marshal, _ := json.Marshal(connector.RoomInfo)
		logger.Infof("[kuaishou.connnector]GetRoomInfo SKIP for ALREADY have value:%s", marshal)
		return connector.RoomInfo
	}

	liveRoomId, liveRoomCaption, err := connector.GetLiveRoomId()
	if err != nil {
		return nil
	}
	connector.GetWebSocketInfo(liveRoomId)
	connector.RoomInfo = &domain.RoomInfo{
		RoomId:    liveRoomId,
		RoomTitle: liveRoomCaption,
	}

	logger.Infof("[kuaishou.connnector][GetRoomInfo]return with roomId:%s roomTitle:%s", liveRoomId, liveRoomCaption)
	return connector.RoomInfo
}

func (connector *ConnectorStrategy) StartListen(localConn *websocket.Conn) {
	logger.Infof("[kuaishou.connnector]StartListen for room:%s", connector.RoomInfo.RoomId)
	connector.IsStart = true
	dataRegistry := data.GetDataRegistry()
	for {
		_, message, err := connector.conn.ReadMessage()

		select {
		case <-connector.stopChan:
			// Stop signal received, exit the goroutine
			logger.Infof("⚓️♪ [kuaishou.ConnectorStrategy] StartListen BREAKED by c.stopChan")
			return
		default:
			if err != nil {
				logger.Errorf("[kuaishou.ConnectorStrategy] StartListen fail with reason: %e", err)
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
	logger.Infof("🎦[kuaishou.ConnectorStrategy] Stop douyin for url: %s title: %s", connector.liveUrl, title)
}

func (connector *ConnectorStrategy) IsAlive() bool {
	return connector.IsStop
}

func (connector *ConnectorStrategy) GetLiveRoomId() (string, string, error) {
	req, err := http.NewRequest("GET", connector.liveUrl, nil)
	if err != nil {
		fmt.Println("[kuaishou.connector][GetLiveRoomId]Failed to create HTTP request:", err)
		return "", "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header = http.Header{}
	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	req.Header.Add("Cookie", HeaderCookieValue)
	req.Header.Add("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Infof("[kuaishou.connector]Failed to send HTTP request:%e", err)
		return "", "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[kuaishou.connector]getRoomInfoByRequest error when read body: %e", err)
		return "", "", err
	}

	regExp := regexp.MustCompile(RoomInfoRegExp)
	bodyString := string(bodyBytes)
	jsonMatches := regExp.FindStringSubmatch(bodyString)
	jsonData := jsonMatches[1]
	logger.Infof("[kuaishou.connector]roomData: %s", jsonData)

	root, err := ajson.Unmarshal([]byte(jsonData))
	liveRoomIdNodes, err := root.JSONPath("$.liveroom.liveStream.id")
	liveRoomId := ""
	for _, node := range liveRoomIdNodes {
		liveRoomId = node.MustString()
		break
	}

	liveCaptionNodes, err := root.JSONPath("$.liveroom.liveStream.caption")
	liveRoomCaption := ""
	for _, node := range liveCaptionNodes {
		liveRoomCaption = node.MustString()
		break
	}

	connector.ExtensionInfo.liveRoomId = liveRoomId
	if liveRoomId == "" {
		return "", "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	return liveRoomId, liveRoomCaption, nil
}

func (connector *ConnectorStrategy) GetWebSocketInfo(liveRoomId string) (*ExtensionInfo, error) {
	requestUrl := fmt.Sprintf(RoomInfoRequestURLPattern, liveRoomId)
	logger.Infof("GetWebSocketInfo with liveRoomId: %s to request url: %s", liveRoomId, requestUrl)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		fmt.Println("[kuaishou.connector]Failed to create HTTP request:", err)
		return nil, fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	req.Header.Add("Cookie", HeaderCookieValue)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Infof("[kuaishou.connector]Failed to send HTTP request:%e", err)
		return nil, fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	logger.Infof("[kuaishou.connector]GetWebSocketInfo with result:%s", bodyString)

	root, err := ajson.Unmarshal([]byte(bodyString))
	tokenNodes, err := root.JSONPath("$.data.token")
	token := ""
	for _, node := range tokenNodes {
		token = node.MustString()
		break
	}

	webSocketUrlNodes, err := root.JSONPath("$.data.websocketUrls[0]")
	webSocketUrl := ""
	for _, node := range webSocketUrlNodes {
		webSocketUrl = node.MustString()
		break
	}

	if token == "" || webSocketUrl == "" {
		logger.Errorf("[kuaishou.connector]GetWebSocketInfo fail to get requestUrl:%s token:%s", webSocketUrl, token)
		return nil, fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	connector.ExtensionInfo.token = token
	connector.ExtensionInfo.webSocketUrl = webSocketUrl
	logger.Infof("[kuaishou.connector]GetWebSocketInfo with requestUrl:%s token:%s", webSocketUrl, token)
	return &connector.ExtensionInfo, nil
}

func (connector *ConnectorStrategy) StartHeartbeat(ws *websocket.Conn) {
	for {
		select {
		case <-connector.stopChan:
			logger.Warnf("💔[kuaishou.connector]StartHeartbeat stop by stopChan notify for roomId:%s", connector.RoomInfo.RoomId)
			return
		default:
			time.Sleep(20 * time.Second)
			obj := &kuaishou_protostub.CSWebHeartbeat{
				PayloadType: 1,
				Payload: &kuaishou_protostub.CSWebHeartbeat_Payload{
					Timestamp: uint64(time.Now().Unix()),
				},
			}

			data, err := proto.Marshal(obj)
			if err != nil {
				logger.Infof("[kuaishou.connector][StartHeartbeat] proto.Marshal fail with err:%e", err)
				continue
			}
			err = ws.WriteMessage(websocket.BinaryMessage, data)
			logger.Infof("❤️[kuaishou.connector][StartHeartbeat] [发送心跳]")

			if err != nil {
				log.Println("[kuaishou.connector][StartHeartbeat] [发送数据失败]", err)
				continue
			}
		}
	}
}

func (connector *ConnectorStrategy) getPageID() string {
	charset := PageIdCharacterSet
	pageID := ""
	for i := 0; i < 16; i++ {
		pageID += string(charset[rand.Intn(len(charset))])
	}
	pageID += "_"
	pageID += fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	return pageID
}
