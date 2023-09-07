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
		logger.Infof("🎦[kuaishou.ConnectorStrategy] Start kuaishou fail for url: %s", connector.liveUrl)
		return constant.INVALID_LIVE_URL
	}

	err := connector.GetWebSocketInfo(roomInfo.RoomId)
	if err != nil {
		log.Fatal(err)
	}

	header := connector.ExtensionInfo.Headers

	conn, dialResp, err := websocket.DefaultDialer.Dial(connector.ExtensionInfo.webSocketUrl, header)
	if err != nil {
		log.Fatal(err)
	}
	marshal, _ := json.Marshal(dialResp)
	logger.Infof("[kuaishou.connnector][connect]dial with response:%s", marshal)
	connector.OnOpen(conn)

	connector.conn = conn
	return constant.SUCCESS
}

func (connector *ConnectorStrategy) GetRoomInfo() *domain.RoomInfo {
	if connector.RoomInfo != nil {
		return connector.RoomInfo
	}
	liveRoomId, err := connector.GetLiveRoomId()
	if err != nil {
		return nil
	}
	connector.GetWebSocketInfo(liveRoomId)
	connector.RoomInfo = &domain.RoomInfo{
		RoomId: liveRoomId,
	}
	return connector.RoomInfo
}

func (connector *ConnectorStrategy) StartListen(localConn *websocket.Conn) {
	logger.Infof("[kuaishou.connnector]StartListen for room:%s", connector.RoomInfo.RoomId)
	connector.IsStart = true
	for {
		_, message, err := connector.conn.ReadMessage()

		//select {
		//case <-connector.stopChan:
		//	// Stop signal received, exit the goroutine
		//	logger.Infof("⚓️♪ [douyin.ConnectorStrategy] StartListen BREAKED by c.stopChan")
		//	return
		//default:
		if err != nil {
			logger.Errorf("[douyin.ConnectorStrategy] StartListen fail with reason: %e", err)
			connector.Stop()
			return
		}
		connector.OnMessage(message, connector.conn, localConn)
		//}
	}

}

func (connector *ConnectorStrategy) Stop() {

}

func (connector *ConnectorStrategy) IsAlive() bool {
	return false
}

func (connector *ConnectorStrategy) GetLiveRoomId() (string, error) {
	req, err := http.NewRequest("GET", connector.liveUrl, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header = http.Header{}
	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	//req.Header.Add("Postman-Token", uuid.NewString())
	req.Header.Add("Cookie", HeaderCookieValue)
	req.Header.Add("Cache-Control", "no-cache")
	//req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	//req.Header.Add("Sec-Ch-Ua", "\"Chromium\";v=\"116\", \"Not)A;Brand\";v=\"24\", \"Google Chrome\";v=\"116\"")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send HTTP request:", err)
		return "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[kuaishou.connector]getRoomInfoByRequest error when read body: %e", err)
		return "", err
	}

	regexp := regexp.MustCompile(RoomInfoRegExp)
	bodyString := string(bodyBytes)
	jsonMatches := regexp.FindStringSubmatch(bodyString)
	jsonData := jsonMatches[1]
	log.Infof("roomData: %s", jsonData)

	root, err := ajson.Unmarshal([]byte(jsonData))
	liveRoomIdNodes, err := root.JSONPath("$.liveroom.liveStream.id")
	liveRoomId := ""
	for _, node := range liveRoomIdNodes {
		liveRoomId = node.MustString()
		break
	}

	connector.ExtensionInfo.liveRoomId = liveRoomId
	if liveRoomId == "" {
		return "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	return liveRoomId, nil
}

func (connector *ConnectorStrategy) GetWebSocketInfo(liveRoomId string) error {
	requestUrl := fmt.Sprintf(RoomInfoRequestURLPattern, liveRoomId)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	req.Header.Add("Cookie", HeaderCookieValue)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send HTTP request:", err)
		return fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	log.Info("GetWebSocketInfo with result:%s", bodyString)

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

	connector.ExtensionInfo.token = token
	connector.ExtensionInfo.webSocketUrl = webSocketUrl
	log.Infof("GetWebSocketInfo with requestUrl:%s token:%s", webSocketUrl, token)
	return nil
}

func (connector *ConnectorStrategy) OnOpen(ws *websocket.Conn) {
	data := connector.connectData()
	log.Println("[kuaishou.connector][onOpen] [建立wss连接]")
	err := ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println("[kuaishou.connector][onOpen] [发送数据失败]", err)
	}

	go connector.KeepHeartBeat(ws)
}

func (connector *ConnectorStrategy) connectData() []byte {
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
		log.Println("[kuaishou.connector][connectData] [序列化失败]", err)
		return nil
	}
	logger.Infof("[kuaishou.connector][connectData] sent data: %s", marshal)
	return marshal
}

func (connector *ConnectorStrategy) heartbeatData() []byte {
	obj := kuaishou_protostub.CSWebHeartbeat{
		PayloadType: 1,
		Payload: &kuaishou_protostub.CSWebHeartbeat_Payload{
			Timestamp: uint64(time.Now().Unix()),
		},
	}

	data, err := json.Marshal(obj)
	if err != nil {
		logger.Infof("[kuaishou.connector][heartbeatData] [序列化失败] :%e", err)
		return nil
	}
	return data
}

func (connector *ConnectorStrategy) KeepHeartBeat(ws *websocket.Conn) {
	for {
		//select {
		//case <-connector.stopChan:
		//	logger.Warnf("[kuaishou.connector]KeepHeartBeat stop by stopChan notify for roomId:%s", connector.RoomInfo.RoomId)
		//	break
		//default:
		time.Sleep(20 * time.Second)
		payload := connector.heartbeatData()
		log.Println("[KeepHeartBeat] [发送心跳]")
		err := ws.WriteMessage(websocket.BinaryMessage, payload)
		if err != nil {
			log.Println("[KeepHeartBeat] [发送数据失败]", err)
		}
		time.Sleep(20 * time.Second)

		//}
	}
}

func (connector *ConnectorStrategy) getPageID() string {
	charset := "-_zyxwvutsrqponmlkjihgfedcba9876543210ZYXWVUTSRQPONMLKJIHGFEDCBA"
	pageID := ""
	for i := 0; i < 16; i++ {
		pageID += string(charset[rand.Intn(len(charset))])
	}
	pageID += "_"
	pageID += fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	return pageID
}
