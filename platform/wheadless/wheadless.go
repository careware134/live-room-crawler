package wheadless

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/util"
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

func NewInstance(Target domain.TargetStruct, stopChan chan struct{}, localConn *websocket.Conn) *HeadlessConnectorStrategy {
	logger.Infof("👓[headless.ConnectorStrategy] NewInstance for url: %s cookie:%s", Target.LiveURL, Target.Cookie)
	return &HeadlessConnectorStrategy{
		Target:    Target,
		stopChan:  stopChan,
		localConn: localConn,
	}
}

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
	logger.Infof("👓[headless.ConnectorStrategy]Connect HEADLESS SKIP for room:%s", c.RoomInfo.RoomTitle)
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
	c.RoomInfo = &domain.RoomInfo{
		RoomTitle: c.Target.LiveURL,
		RoomId:    c.Target.LiveURL,
	}
	return c.RoomInfo
	//req, err := http.NewRequest("GET", c.Target.LiveURL, nil)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	//
	//// construct header to simulate connection
	//req.Header = http.Header{
	//	"Accept":          []string{"*/*"},
	//	"User-Agent":      []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"},
	//	"Accept-Encoding": []string{"gzip, deflate, br"},
	//	"Connection":      []string{"keep-alive"},
	//}
	//
	//// Create a new HTTP connection and send the request
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//defer resp.Body.Close()
	//
	//// Parse the HTML response
	//var bodyReader *bytes.Reader
	//switch resp.Header.Get("Content-Encoding") {
	//case "gzip":
	//	gzipReader, err := gzip.NewReader(resp.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer gzipReader.Close()
	//	body, err := ioutil.ReadAll(gzipReader)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	bodyReader = bytes.NewReader(body)
	//
	//case "deflate":
	//	flateReader := flate.NewReader(resp.Body)
	//	defer flateReader.Close()
	//	body, err := ioutil.ReadAll(flateReader)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	bodyReader = bytes.NewReader(body)
	//
	//default:
	//	body, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	bodyReader = bytes.NewReader(body)
	//}
	//
	//// Parse the HTML response using goquery
	//doc, err := goquery.NewDocumentFromReader(bodyReader)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Get the innerHTML from the <title> tag
	//title := strings.TrimSpace(doc.Text())
	//logger.Infof("head: %s", title)
	//
	//// Get the value of the 'content' attribute from the <meta> tag
	//cid, exists := doc.Find("html head meta[name='lx:cid']").Attr("content")
	//if !exists {
	//	return nil
	//}
	//
	//return &domain.RoomInfo{
	//	RoomTitle: title,
	//	RoomId:    cid,
	//}
}
