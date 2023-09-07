package kuaishou

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spyzhov/ajson"
	"google.golang.org/protobuf/proto"
	"io"
	"live-room-crawler/constant"
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
	token        string
	cookie       string
	webSocketUrl string
	liveRoomId   string
	liveUrl      string
	Headers      http.Header
}

func (t *ConnectorStrategy) Init(liveUrl string) {
	t.liveUrl = liveUrl

}

func (t *ConnectorStrategy) GetLiveRoomId() (string, error) {
	req, err := http.NewRequest("GET", t.liveUrl, nil)
	if err != nil {
		fmt.Println("Failed to create HTTP request:", err)
		return "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header = http.Header{}
	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	req.Header.Add("Postman-Token", uuid.NewString())
	req.Header.Add("Cookie", HeaderCookieValue)
	req.Header.Add("Cache-Control", "no-cache")

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

	t.liveRoomId = liveRoomId
	if liveRoomId == "" {
		return "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	return liveRoomId, nil
}

func (t *ConnectorStrategy) GetWebSocketInfo(liveRoomId string) error {
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

	t.token = token
	t.webSocketUrl = webSocketUrl
	log.Infof("GetWebSocketInfo with requestUrl:%s token:%s", webSocketUrl, token)
	return nil
}

func (t *ConnectorStrategy) WssServerStart() {
	rid, err := t.GetLiveRoomId()
	if err != nil {
		log.Fatal(err)
	}
	err = t.GetWebSocketInfo(rid)
	if err != nil {
		log.Fatal(err)
	}

	c, _, err := websocket.DefaultDialer.Dial(t.webSocketUrl, nil)
	t.OnOpen(c)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("fail to read message with error:%e", err)
		}

		onMessage(message)
	}
}

func onMessage(message []byte) {
	wssPackage := &kuaishou_protostub.SocketMessage{}
	if err := proto.Unmarshal(message, wssPackage); err != nil {
		log.Printf("[onMessage] [无法解析的数据包⚠️] %v\n", err)
		return
	}

	switch wssPackage.PayloadType {
	case kuaishou_protostub.PayloadType_SC_ENTER_ROOM_ACK:
		parseEnterRoomAckPack(wssPackage.Payload)
	case kuaishou_protostub.PayloadType_SC_HEARTBEAT_ACK:
		parseHeartBeatPack(wssPackage.Payload)
	case kuaishou_protostub.PayloadType_SC_FEED_PUSH:
		parseFeedPushPack(wssPackage.Payload)
	case kuaishou_protostub.PayloadType_SC_LIVE_WATCHING_LIST:
		parseSCWebLiveWatchingUsers(wssPackage.Payload)
	default:
		jsonData, err := json.Marshal(wssPackage)
		if err != nil {
			log.Printf("[onMessage] [无法解析的数据包⚠️] %v\n", err)
			return
		}
		log.Printf("[onMessage] [无法解析的数据包⚠️] %s\n", jsonData)
	}
}

func (c *ConnectorStrategy) OnOpen(ws *websocket.Conn) {
	data := c.connectData()
	log.Println("[onOpen] [建立wss连接]")
	err := ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Println("[onOpen] [发送数据失败]", err)
	}

	go c.keepHeartBeat(ws)
}

func (c *ConnectorStrategy) connectData() []byte {
	obj := kuaishou_protostub.CSWebEnterRoom{
		PayloadType: 200,
		Payload: &kuaishou_protostub.CSWebEnterRoom_Payload{
			Token:        c.token,
			LiveStreamId: c.liveRoomId,
			PageId:       c.getPageID(),
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("[connectData] [序列化失败]", err)
		return nil
	}
	return data
}

func (c *ConnectorStrategy) heartbeatData() []byte {
	obj := kuaishou_protostub.CSWebHeartbeat{
		PayloadType: 1,
		Payload: &kuaishou_protostub.CSWebHeartbeat_Payload{
			Timestamp: uint64(time.Now().Unix()),
		},
	}

	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("[heartbeatData] [序列化失败]", err)
		return nil
	}
	return data
}

func (c *ConnectorStrategy) keepHeartBeat(ws *websocket.Conn) {
	for {
		time.Sleep(20 * time.Second)
		payload := c.heartbeatData()
		log.Println("[keepHeartBeat] [发送心跳]")
		err := ws.WriteMessage(websocket.BinaryMessage, payload)
		if err != nil {
			log.Println("[keepHeartBeat] [发送数据失败]", err)
		}
	}
}

func (c *ConnectorStrategy) getPageID() string {
	charset := "-_zyxwvutsrqponmlkjihgfedcba9876543210ZYXWVUTSRQPONMLKJIHGFEDCBA"
	pageID := ""
	for i := 0; i < 16; i++ {
		pageID += string(charset[rand.Intn(len(charset))])
	}
	pageID += "_"
	pageID += fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	return pageID
}
