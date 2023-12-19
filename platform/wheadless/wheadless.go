package wheadless

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/util"
	"net/http"
	"strings"
)

var (
	log = util.Logger()
)

type HeadlessConnectorStrategy struct {
	Target    domain.TargetStruct
	RoomInfo  *domain.RoomInfo
	conn      *websocket.Conn
	localConn *websocket.Conn
	IsStart   bool
	IsStop    bool
	stopChan  chan struct{}
}

var (
	logger = util.Logger()
)

func (c *HeadlessConnectorStrategy) SetRoomInfo(info domain.RoomInfo) {
	marshal, _ := json.Marshal(info)
	c.RoomInfo = &info
	logger.Infof("👓[headless.ConnectorStrategy] SetRoomInfo with value: %s", marshal)
}

func (c *HeadlessConnectorStrategy) VerifyTarget() *domain.CommandResponse {
	info := c.GetRoomInfo()
	responseStatus := constant.SUCCESS
	if info == nil {
		responseStatus = constant.CONNECTION_FAIL
		return &domain.CommandResponse{
			ResponseStatus: responseStatus,
		}
	}

	return &domain.CommandResponse{
		Room:           *info,
		ResponseStatus: responseStatus,
	}
}

func (c *HeadlessConnectorStrategy) Connect() constant.ResponseStatus {
	roomInfo := c.GetRoomInfo()
	if roomInfo == nil {
		logger.Infof("👓[headless.ConnectorStrategy] Connect douyin fail for url: %s", c.Target.LiveURL)
		return constant.INVALID_LIVE_URL
	}

	return constant.SUCCESS
}

func (c *HeadlessConnectorStrategy) StartListen(localConn *websocket.Conn) {
	logger.Infof("👓[headless.ConnectorStrategy]StartListen HEADLESS SKIP for room:%s", c.RoomInfo.RoomTitle)
}

func (c *HeadlessConnectorStrategy) Stop() {
	c.IsStart = false
	c.IsStop = true
	if c.conn != nil {
		c.conn.Close()
	}
	title := ""
	if c.RoomInfo != nil {
		title = c.RoomInfo.RoomTitle
	}
	logger.Infof("👓[headless.ConnectorStrategy]stop for room:%s", title)
}

func (c *HeadlessConnectorStrategy) IsDead() bool {
	return c.IsStop
}

func (c *HeadlessConnectorStrategy) GetRoomInfo() *domain.RoomInfo {
	url := c.Target.LiveURL
	// Send HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the HTML response
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil
	}

	// Get the innerHTML from the <title> tag
	title := strings.TrimSpace(doc.Find("title").Text())

	// Get the value of the 'content' attribute from the <meta> tag
	cid, exists := doc.Find("meta[name='lx:cid']").Attr("content")
	if !exists {
		return nil
	}

	return &domain.RoomInfo{
		RoomTitle: title,
		RoomId:    cid,
	}
}
