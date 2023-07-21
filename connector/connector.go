package connector

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spyzhov/ajson"
	"io"
	"live-room-crawler/listener"
	"live-room-crawler/util"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	roomInfo RoomInfo
	logger   *log.Entry = util.Logger()
)

type RoomInfo struct {
	RoomId    string
	RoomTitle string
	ttwid     string
}

func WssServerStart(roomInfo *RoomInfo) {
	websocketUrl := strings.ReplaceAll(util.WebSocketTemplateURL, util.RoomIdPlaceHolder, roomInfo.RoomId)
	header := http.Header{
		"cookie":     []string{"ttwid=" + roomInfo.ttwid},
		"User-Agent": []string{util.SimulateUserAgent},
	}

	conn, _, err := websocket.DefaultDialer.Dial(websocketUrl, header)

	if err != nil {
		log.Fatalf("fatal to dial websocket! url: %s", websocketUrl, err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		listener.OnMessage(message, conn)
	}
}

func RetrieveRoomInfoFromHttpCall(liveUrl string) *RoomInfo {
	// request room liveUrl
	req, err := http.NewRequest("GET", liveUrl, nil)
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
	liveRoomIdNodes, err := root.JSONPath("$.app.initialState.roomStore.roomInfo.roomId")
	if err != nil {
		panic(err)
	}

	liveRoomId, liveRoomTitle, ttwid := "", "", ""
	for _, node := range liveRoomIdNodes {
		liveRoomId = node.MustString()
		break
	}

	liveRoomTitleNodes, err := root.JSONPath("$.app.initialState.roomStore.roomInfo.room.title")
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

	log.Infof("🎥start to crawl for RoomId: %s title: %s ttwid: %s", liveRoomId, liveRoomTitle, ttwid)
	roomInfo = RoomInfo{
		RoomId:    liveRoomId,
		RoomTitle: liveRoomTitle,
		ttwid:     ttwid,
	}
	return &roomInfo
}
