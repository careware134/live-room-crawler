package kuaishou

import (
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

func (connector *ConnectorStrategy) GetLiveRoomId() (string, string, error) {
	url := connector.Target.LiveURL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("[kuaishou.connector][GetLiveRoomId]Failed to create HTTP request:", err)
		return "", "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}

	req.Header = http.Header{}
	req.Header.Add("Accept", HeaderAcceptValue)
	req.Header.Add("User-Agent", HeaderAgentValue)
	cookie := connector.pickupCookie()
	req.Header.Add("cookie", cookie)
	util.Logger().Infof("[kuaishou.connector]request to get roomData url:%s cookie: %s", url, cookie)
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
	if jsonMatches == nil {
		util.Logger().Infof("[kuaishou.connector]GetRoomId fail mostly its not a leggal url: %s", url)
		return "", "", fmt.Errorf(constant.INVALID_LIVE_URL.Code)
	}
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

func (connector *ConnectorStrategy) GetWebSocketInfo(liveRoomId string) (*domain.RoomInfo, error) {
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

	req.Header.Add("cookie", cookie)
	util.Logger().Infof("[kuaishou.connector]request GetWebSocketInfo url:%s cookie: %s", requestUrl, cookie)
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

	connector.RoomInfo.Token = token
	connector.RoomInfo.WebSocketUrl = webSocketUrl
	util.Logger().Infof("[kuaishou.connector]GetWebSocketInfo with requestUrl:%s token:%s", webSocketUrl, token)
	return connector.RoomInfo, nil
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

func (connector *ConnectorStrategy) pickupCookie2() string {
	return HeaderCookieValue7
}

func (connector *ConnectorStrategy) pickupCookie() string {
	cookie := connector.Target.Cookie
	if cookie != "" {
		util.Logger().Infof("[kauishou.connector]pickupCookie from request: %s", cookie)
		return cookie
	}

	if connector.ExtensionInfo.cookie != "" {
		cookie = connector.ExtensionInfo.cookie
		util.Logger().Infof("[kauishou.connector]pickupCookie by RESURE extensionInfo field, cookie: %s", cookie)
		return cookie
	}

	cookieList := data.GetDataRegistry().GetCookieList(connector.localConn, string(domain.KUAISHOU))
	length := len(cookieList)
	if cookieList != nil && length > 0 {
		randomIndex := rand.Intn(length)
		cookie = cookieList[randomIndex]
		connector.ExtensionInfo.cookie = cookie
		util.Logger().Infof("[kauishou.connector]pickupCookie from hornor list: %s", cookie)
		return cookie
	}

	util.Logger().Infof("[kauishou.connector]pickupCookie by default: %s", cookie)
	return HeaderCookieValue4
}
