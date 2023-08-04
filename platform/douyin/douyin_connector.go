package douyin

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spyzhov/ajson"
	"io"
	"live-room-crawler/common"
	"live-room-crawler/constant"
	"live-room-crawler/util"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	logger = util.Logger()
)

type ConnectorStrategy struct {
	liveUrl  string
	RoomInfo *common.RoomInfo
	conn     *websocket.Conn
}

func NewInstance(liveUrl string) *ConnectorStrategy {
	logger.Infof("♪ NewInstance for url: %s", liveUrl)
	return &ConnectorStrategy{
		liveUrl: liveUrl,
	}
}

func (c *ConnectorStrategy) Start() constant.ResponseStatus {
	roomInfo := c.GetRoomInfo()
	if roomInfo == nil {
		logger.Infof("♪ Start douyin fail for url: %s title: %s", c.liveUrl, c.RoomInfo.RoomTitle)
		return constant.INVALID_LIVE_URL
	}

	websocketUrl := strings.ReplaceAll(util.WebSocketTemplateURL, util.RoomIdPlaceHolder, roomInfo.RoomId)
	header := http.Header{
		"cookie":     []string{"ttwid=" + roomInfo.Ttwid},
		"User-Agent": []string{util.SimulateUserAgent},
	}

	conn, _, err := websocket.DefaultDialer.Dial(websocketUrl, header)
	c.conn = conn
	logger.Infof("♪ Start douyin success for url: %s title: %s", c.liveUrl, c.RoomInfo.RoomTitle)

	if err != nil {
		logger.Fatalf("fatal to dial websocket! url: %s error:%e", websocketUrl, err)
		return constant.CONNECTION_FAIL
	}

	return constant.SUCCESS
}

func (c *ConnectorStrategy) Stop() {
	c.conn.Close()
	logger.Infof("♪ Stop douyin for url: %s titl: %s", c.liveUrl, c.RoomInfo.RoomTitle)
}

func (c *ConnectorStrategy) GetRoomInfo() *common.RoomInfo {
	// request room liveUrl
	req, err := http.NewRequest("GET", c.liveUrl, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// construct header to simulate client
	req.Header = http.Header{
		"Accept":     []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"},
		"Cookie":     []string{"__ac_nonce=0638733a400869171be51"},
	}

	// Create a new HTTP client and send the request
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
	renderDataRegex := regexp.MustCompile(`<script id="RENDER_DATA" type="application/json">(.*?)</script>`)

	regexGroup := renderDataRegex.FindStringSubmatch(body)
	if len(regexGroup) < 2 {
		fmt.Println("No render data found")
		return nil
	}
	renderData := regexGroup[1]
	renderData, err = url.QueryUnescape(renderData) // liveUrl.QueryUnescape(renderData)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	root, err := ajson.Unmarshal([]byte(renderData))
	if err != nil {
		panic(err)
	}

	// Retrieve the value of the "name" field using JSONPath
	liveRoomIdNodes, err := root.JSONPath("$.app.initialState.roomStore.RoomInfo.roomId")
	if err != nil {
		panic(err)
	}

	liveRoomId, liveRoomTitle, ttwid := "", "", ""
	for _, node := range liveRoomIdNodes {
		liveRoomId = node.MustString()
		break
	}

	liveRoomTitleNodes, err := root.JSONPath("$.app.initialState.roomStore.RoomInfo.room.title")
	if err != nil {
		panic(err)
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

	logger.Infof("🎥start to crawl for RoomId: %s title: %s ttwid: %s", liveRoomId, liveRoomTitle, ttwid)
	c.RoomInfo = &common.RoomInfo{
		RoomId:    liveRoomId,
		RoomTitle: liveRoomTitle,
		Ttwid:     ttwid,
	}

	logger.Infof("♪ GetRoomInfo douyin for url: %s titlt: %s", c.liveUrl, c.RoomInfo.RoomTitle)
	return c.RoomInfo
}
