package douyin

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/util"
	"net/http"
	"regexp"
	"strings"
)

var (
	logger = util.Logger()
)

type ConnectorStrategy struct {
	liveUrl   string
	RoomInfo  *domain.RoomInfo
	conn      *websocket.Conn
	localConn *websocket.Conn
	IsStart   bool
	IsStop    bool
	stopChan  chan struct{}
}

func NewInstance(liveUrl string, stopChan chan struct{}) *ConnectorStrategy {
	logger.Infof("♪ [douyin.ConnectorStrategy] NewInstance for url: %s", liveUrl)
	return &ConnectorStrategy{
		liveUrl:  liveUrl,
		stopChan: stopChan,
	}
}

func (c *ConnectorStrategy) Connect(localConn *websocket.Conn) constant.ResponseStatus {
	roomInfo := c.GetRoomInfo()
	if roomInfo == nil {
		logger.Infof("♪ [douyin.ConnectorStrategy] Start douyin fail for url: %s", c.liveUrl)
		return constant.INVALID_LIVE_URL
	}

	websocketUrl := strings.ReplaceAll(WebSocketTemplateURL, RoomIdPlaceHolder, roomInfo.RoomId)
	header := http.Header{
		"cookie":     []string{"ttwid=" + roomInfo.Ttwid},
		"User-Agent": []string{SimulateUserAgent},
	}

	conn, _, err := websocket.DefaultDialer.Dial(websocketUrl, header)
	c.conn = conn
	logger.Infof("♪ [douyin.ConnectorStrategy] Start douyin success for url: %s title: %s", c.liveUrl, c.RoomInfo.RoomTitle)

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
	logger.Infof("♪[douyin.ConnectorStrategy] Stop douyin for url: %s title: %s", c.liveUrl, title)
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
	req, err := http.NewRequest("GET", c.liveUrl, nil)
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
		fmt.Println(err)
		return nil
	}

	body := string(bodyBytes)
	liveRoomId, liveRoomTitle, ttwid := "", "", ""

	regex := regexp.MustCompile(`roomId\\":\\"(\d+)\\"`)
	match := regex.FindStringSubmatch(body)
	if len(match) > 1 {
		liveRoomId = match[1]
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
		RoomId:    liveRoomId,
		RoomTitle: liveRoomTitle,
		Ttwid:     ttwid,
	}

	return c.RoomInfo
}
