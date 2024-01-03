package douyin

import (
	"encoding/json"
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

func NewInstance(Target domain.TargetStruct, stopChan chan struct{}, localConn *websocket.Conn) *ConnectorStrategy {
	logger.Infof("♪ [douyin.ConnectorStrategy] NewInstance for url: %s cookie:%s", Target.LiveURL, Target.Cookie)
	return &ConnectorStrategy{
		Target:    Target,
		stopChan:  stopChan,
		localConn: localConn,
	}
}

func (c *ConnectorStrategy) SetRoomInfo(info domain.RoomInfo) {
	marshal, _ := json.Marshal(info)
	c.RoomInfo = &info
	logger.Infof("♪ [douyin.ConnectorStrategy] SetRoomInfo with value: %s", marshal)
}

func (c *ConnectorStrategy) VerifyTarget() *domain.CommandResponse {
	info := c.GetRoomInfo()
	responseStatus := constant.SUCCESS
	if info == nil {
		responseStatus = constant.CONNECTION_FAIL
		return &domain.CommandResponse{
			ResponseStatus: responseStatus,
		}
	}

	return &domain.CommandResponse{
		Room:           info,
		ResponseStatus: responseStatus,
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
	logger.Infof("♪[douyin.ConnectorStrategy] Stop kuaishou for url: %s title: %s", c.Target.LiveURL, title)
}

func (c *ConnectorStrategy) IsDead() bool {
	return c.IsStop
}

func (c *ConnectorStrategy) GetRoomInfo() *domain.RoomInfo {
	if c.RoomInfo != nil {
		logger.Infof("♪[douyin.ConnectorStrategy] GetRoomInfo return directly!")
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
		"Cookie":     []string{"ttwid=1%7CcMckL3_tA47rt11Ektt038lvR_qO4_w9TPINxgIVWBQ%7C1704274070%7Cdceb110148d3e384753ac64e55ea92efa1389600ecb71c010aba193ed82dfb19; has_avx2=null; device_web_cpu_core=12; device_web_memory_size=8; live_use_vvc=%22false%22; xgplayer_user_id=136259113385; csrf_session_id=90bccab7b3988b8c03b8e775b678414b; ttcid=4cfd498bef754336baea64cd3476676831; FORCE_LOGIN=%7B%22videoConsumedRemainSeconds%22%3A180%7D; __ac_nonce=065952d1b001f10f86f0b; __ac_signature=_02B4Z6wo00f014jIhpAAAIDCBclR1OcXzYuI6IIAAIevjD9PRMafTVjOd0.JDui.ndFCusoT6dMpOIzZIGi0vgkWzP0CSXfH0nkKdBTmyuR-GMgDaQQN8xkan8qdNwkV9f21uhxOFK07u63r41; webcast_local_quality=sd; webcast_leading_last_show_time=1704276190454; webcast_leading_total_show_times=2; pwa2=%220%7C0%7C3%7C0%22; download_guide=%223%2F20240103%2F0%22; xg_device_score=7.802204888412783; msToken=8VpuTPur6-0ShNqAnvMsPVT6NsiuzNgjKTKrYHbBKMQVsT0Sp5WZtC6-fQDUxSDlpuRiBzj7UUSJPc6p3Wq4Or3UqdKKkJEHIVErONlhZKmyuL4vWaTKaPWF4NeM; tt_scid=Sy4boUfYrwdXqtJ917AF57iKJhzV1yQS-5KNpTlP6irjDGwMkVGIy8bcEcRE33Qx32ff; __live_version__=%221.1.1.6845%22; live_can_add_dy_2_desktop=%221%22; msToken=fHCLfxKr4SqfHl8kSyR4jew35LV-niUiTNJqmZUDkdu74ZGueJ0a6b6opVP8Rhjkl0XEziqAnTgFBOKJ4fKJ67sFSPRf1-r0VA6kC3ty-v35WYKP-gqnNNvfe8LL; IsDouyinActive=false"},
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
