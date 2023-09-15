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
	liveUrl      string
	token        string
	cookie       string
	webSocketUrl string
	liveRoomId   string
	Headers      http.Header
}

var (
	logger = util.Logger()
)

func NewInstance(Target domain.TargetStruct, stopChan chan struct{}) *ConnectorStrategy {
	util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] NewInstance for url: %s cookie:%s", Target.LiveURL, Target.Cookie)
	return &ConnectorStrategy{
		Target:   Target,
		stopChan: stopChan,
	}
}

func (connector *ConnectorStrategy) Connect() constant.ResponseStatus {
	roomInfo := connector.GetRoomInfo()
	if roomInfo == nil {
		util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] fail to get kuaishou roomInfo with url: %s", connector.Target.LiveURL)
		return constant.FAIL_GETTING_ROOM_INFO
	}

	_, err := connector.GetWebSocketInfo(roomInfo.RoomId)
	if err != nil {
		return constant.FAIL_GETTING_SOCKET_INFO
	}

	header := connector.ExtensionInfo.Headers
	url := connector.ExtensionInfo.webSocketUrl
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
			Token:        connector.ExtensionInfo.token,
			LiveStreamId: connector.ExtensionInfo.liveRoomId,
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
		return nil
	}
	connector.GetWebSocketInfo(liveRoomId)
	connector.RoomInfo = &domain.RoomInfo{
		RoomId:    liveRoomId,
		RoomTitle: liveRoomCaption,
	}

	util.Logger().Infof("[kuaishou.connnector][GetRoomInfo]return with roomId:%s roomTitle:%s", liveRoomId, liveRoomCaption)
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
	util.Logger().Infof("🎦[kuaishou.ConnectorStrategy] Stop douyin for url: %s title: %s", connector.Target.LiveURL, title)
}

func (connector *ConnectorStrategy) IsAlive() bool {
	return connector.IsStop
}

func (connector *ConnectorStrategy) GetLiveRoomId() (string, string, error) {
	req, err := http.NewRequest("GET", connector.Target.LiveURL, nil)
	if err != nil {
		fmt.Println("[kuaishou.connector][GetLiveRoomId]Failed to create HTTP request:", err)
		return "", "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header = http.Header{}
	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	cookie := connector.pickupCookie()
	req.Header.Add("Cookie", cookie)
	util.Logger().Infof("[kuaishou.connector]reques roomData url:%s cookie: %s", connector.Target.LiveURL, cookie)
	//req.Header.Add("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.Logger().Infof("[kuaishou.connector]Failed to send HTTP request:%e", err)
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
	util.Logger().Infof("[kuaishou.connector]roomData: %s", jsonData)

	root, err := ajson.Unmarshal([]byte(jsonData))
	liveRoomIdNodes, err := root.JSONPath("$.liveroom..liveStream.id")
	//liveRoomIdNodes, err := root.JSONPath("$.liveroom.playList[0].liveStream.id")
	liveRoomId := ""
	for _, node := range liveRoomIdNodes {
		liveRoomId = node.MustString()
		break
	}

	liveCaptionNodes, err := root.JSONPath("$.liveroom..liveStream.caption")
	//liveCaptionNodes, err := root.JSONPath("$.liveroom.playList[0].liveStream.caption")

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
	util.Logger().Infof("GetWebSocketInfo with liveRoomId: %s to request url: %s", liveRoomId, requestUrl)

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		fmt.Println("[kuaishou.connector]Failed to create HTTP request:", err)
		return nil, fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	cookie := connector.pickupCookie()

	req.Header.Add("Cookie", cookie)
	util.Logger().Infof("[kuaishou.connector]request GetWebSocketInfo url:%s cookie: %s", requestUrl, cookie)
	req.Header.Add("Cookie", cookie)
	//req.Header.Add("Cookie", HeaderCookieValue2)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		util.Logger().Infof("[kuaishou.connector]Failed to send HTTP request:%e", err)
		return nil, fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	util.Logger().Infof("[kuaishou.connector]GetWebSocketInfo with result:%s", bodyString)

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
		util.Logger().Errorf("[kuaishou.connector]GetWebSocketInfo fail to get requestUrl:%s token:%s", webSocketUrl, token)
		return nil, fmt.Errorf(constant.FAIL_GETTING_SOCKET_INFO.Code)
	}

	connector.ExtensionInfo.token = token
	connector.ExtensionInfo.webSocketUrl = webSocketUrl
	util.Logger().Infof("[kuaishou.connector]GetWebSocketInfo with requestUrl:%s token:%s", webSocketUrl, token)
	return &connector.ExtensionInfo, nil
}

func (connector *ConnectorStrategy) StartHeartbeat(ws *websocket.Conn) {
	for {
		select {
		case <-connector.stopChan:
			util.Logger().Warnf("💔[kuaishou.connector]StartHeartbeat stop by stopChan notify for roomId:%s", connector.RoomInfo.RoomId)
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
				util.Logger().Infof("[kuaishou.connector][StartHeartbeat] proto.Marshal fail with err:%e", err)
				continue
			}
			err = ws.WriteMessage(websocket.BinaryMessage, data)
			util.Logger().Infof("❤️[kuaishou.connector][StartHeartbeat] [发送心跳]")

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

func (connector *ConnectorStrategy) pickupCookie() string {
	cookie := connector.Target.Cookie
	if cookie != "" {
		util.Logger().Infof("[kauishou.connector]pickupCookie from request: %s", cookie)
		return cookie
	}

	cookieList := data.GetDataRegistry().GetCookieList(connector.localConn, "kuaishou")
	length := len(cookieList)
	if cookieList != nil && length > 0 {
		randomIndex := rand.Intn(length)
		cookie = cookieList[randomIndex]
		util.Logger().Infof("[kauishou.connector]pickupCookie from hornor list: %s", cookie)
		return cookie
	}

	util.Logger().Infof("[kauishou.connector]pickupCookie by default: %s", cookie)
	return HeaderCookieValue4
}
