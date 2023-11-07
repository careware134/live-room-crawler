package douyin

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spyzhov/ajson"
	"io"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/util"
	"net/http"
	"strings"
)

var (
	logger = util.Logger()
)

type ConnectorStrategy struct {
	Target    domain.TargetStruct
	RoomInfo  *domain.RoomInfo
	conn      *websocket.Conn
	localConn *websocket.Conn
	IsStart   bool
	IsStop    bool
	stopChan  chan struct{}
}

func (c *ConnectorStrategy) VerifyTarget() *domain.CommandResponse {
	info := c.GetRoomInfo()
	responseStatus := constant.SUCCESS
	if info == nil {
		responseStatus = constant.CONNECTION_FAIL
	}

	return &domain.CommandResponse{
		Room:           *info,
		ResponseStatus: responseStatus,
	}
}

func NewInstance(Target domain.TargetStruct, stopChan chan struct{}, localConn *websocket.Conn) *ConnectorStrategy {
	logger.Infof("♪ [douyin.ConnectorStrategy] NewInstance for url: %s cookie:%s", Target.LiveURL, Target.Cookie)
	return &ConnectorStrategy{
		Target:    Target,
		stopChan:  stopChan,
		localConn: localConn,
	}
}

func (c *ConnectorStrategy) Connect() constant.ResponseStatus {
	roomInfo := c.GetRoomInfo()
	if roomInfo == nil {
		logger.Infof("♪ [douyin.ConnectorStrategy] Start douyin fail for url: %s", c.Target.LiveURL)
		return constant.INVALID_LIVE_URL
	}

	websocketUrl := strings.ReplaceAll(WebSocketTemplateURL, RoomIdPlaceHolder, roomInfo.RoomId)
	header := http.Header{
		"cookie":     []string{"ttwid=" + roomInfo.Token},
		"User-Agent": []string{SimulateUserAgent},
	}

	conn, _, err := websocket.DefaultDialer.Dial(websocketUrl, header)
	c.conn = conn
	logger.Infof("♪ [douyin.ConnectorStrategy] Start douyin success for url: %s title: %s", c.Target.LiveURL, c.RoomInfo.RoomTitle)

	if err != nil {
		logger.Errorf(" [douyin.ConnectorStrategy] fail to dial websocket! url: %s error:%e", websocketUrl, err)
		return constant.CONNECTION_FAIL
	}
	return constant.SUCCESS
}

func (c *ConnectorStrategy) StartListen(localConn *websocket.Conn) {
	logger.Infof("StartListen for room:%s", c.RoomInfo.RoomTitle)
	c.IsStart = true
	for {
		_, message, err := c.conn.ReadMessage()

		select {
		case <-c.stopChan:
			// Stop signal received, exit the goroutine
			logger.Infof("⚓️♪ [douyin.ConnectorStrategy] StartListen BREAKED by c.stopChan")
			return
		default:
			if err != nil {
				logger.Errorf("[douyin.ConnectorStrategy] StartListen fail with reason: %e", err)
				c.Stop()
				return
			}
			c.OnMessage(message, c.conn, localConn)
		}
	}
}

func (c *ConnectorStrategy) Stop() {
	c.IsStart = false
	c.IsStop = true
	if c.conn != nil {
		c.conn.Close()
	}
	title := ""
	if c.RoomInfo != nil {
		title = c.RoomInfo.RoomTitle
	}
	logger.Infof("♪[kuaishou.ConnectorStrategy] Stop kuaishou for url: %s title: %s", c.Target.LiveURL, title)
}

func (c *ConnectorStrategy) IsAlive() bool {
	return c.IsStop
}

func (c *ConnectorStrategy) GetRoomInfo() *domain.RoomInfo {
	if c.RoomInfo != nil {
		return c.RoomInfo
	}
	return c.getRoomInfoByRequest()
}

func (c *ConnectorStrategy) getRoomInfoByRequest() *domain.RoomInfo {
	// request room liveUrl
	req, err := http.NewRequest("GET", c.Target.LiveURL, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// construct header to simulate connection
	req.Header = http.Header{
		"Accept":     []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"},
		"Cookie":     []string{"__ac_nonce=0638733a400869171be51"},
	}

	// Create a new HTTP connection and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("[douyin.connector]getRoomInfoByRequest error when read body: %e", err)
		return nil
	}

	body := string(bodyBytes)
	liveRoomId, liveRoomTitle, ttwid := "", "", ""

	//regex := regexp.MustCompile(`roomId\\":\\"(\d+)\\"`)
	//match := regex.FindStringSubmatch(body)
	//if len(match) > 1 {
	//	liveRoomId = match[1]
	//}

	startIndex := strings.LastIndex(body, RoomInfoJsonStartTag1)
	if startIndex < 0 {
		startIndex = strings.LastIndex(body, RoomInfoJsonStartTag2)
	}
	if startIndex < 0 {
		logger.Errorf("[douyin.connector]getRoomInfoByRequest fail to find roomInfo json: %e", err)
		return nil
	}

	jsonNotEnding := body[startIndex:]
	startIndex = strings.Index(jsonNotEnding, "{")
	roomInfoJson := util.FindJsonString(jsonNotEnding, startIndex)

	root, err := ajson.Unmarshal([]byte(roomInfoJson))
	if err != nil {
		logger.Errorf("[douyin.connector]getRoomInfoByRequest error when Unmarshal Json: %e", err)
		return nil
	}

	// Retrieve the value of the "name" field using JSONPath
	liveRoomIdNodes, err := root.JSONPath("$.room.id_str")
	if err != nil {
		logger.Errorf("[douyin.connector]getRoomInfoByRequest error when find liveRoomId: %e", err)
		return nil
	}

	for _, node := range liveRoomIdNodes {
		liveRoomId = node.MustString()
		break
	}

	liveRoomTitleNodes, err := root.JSONPath("$.room.title")
	if err != nil {
		logger.Warnf("[douyin.connector]getRoomInfoByRequest error when find liveRoomTitleNodes: %e", err)
	}
	for _, node := range liveRoomTitleNodes {
		liveRoomTitle = node.MustString()
		break
	}

	data := resp.Cookies()
	for _, cookie := range data {
		if cookie.Name == "ttwid" {
			ttwid = cookie.Value
			break
		}
	}

	logger.Infof("♪  [douyin.ConnectorStrategy] GetRoomInfo for RoomId: %s title: %s ttwid: %s", liveRoomId, liveRoomTitle, ttwid)
	c.RoomInfo = &domain.RoomInfo{
		RoomId:       liveRoomId,
		RoomTitle:    liveRoomTitle,
		Token:        ttwid,
		WebSocketUrl: strings.ReplaceAll(WebSocketTemplateURL, RoomIdPlaceHolder, liveRoomId),
	}

	return c.RoomInfo
}
