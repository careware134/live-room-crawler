package kuaishou

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"live-room-crawler/platform/kuaishou/kuaishou_protostub"
)

func (connector *ConnectorStrategy) OnMessage(message []byte, conn *websocket.Conn, localConn *websocket.Conn) {

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

func parseEnterRoomAckPack(message []byte) {
	scWebEnterRoomAck := &kuaishou_protostub.SCWebEnterRoomAck{}
	if err := proto.Unmarshal(message, scWebEnterRoomAck); err != nil {
		log.Printf("[parseEnterRoomAckPack] [进入房间成功ACK应答👌] fail unmarshal proto: %v", err)
		return
	}
	jsonData, err := json.Marshal(scWebEnterRoomAck)
	if err != nil {
		log.Printf("[parseEnterRoomAckPack] [进入房间成功ACK应答👌]fail unmarshal json: %v", err)
		return
	}
	log.Printf("[parseEnterRoomAckPack] [进入房间成功ACK应答👌] success: %s\n", jsonData)
}

func parseSCWebLiveWatchingUsers(message []byte) {
	scWebLiveWatchingUsers := &kuaishou_protostub.SCWebLiveWatchingUsers{}
	if err := proto.Unmarshal(message, scWebLiveWatchingUsers); err != nil {
		log.Printf("[parseSCWebLiveWatchingUsers] [不知道是啥的数据包🤷] %v\n", err)
		return
	}
	jsonData, err := json.Marshal(scWebLiveWatchingUsers)
	if err != nil {
		log.Printf("[parseSCWebLiveWatchingUsers] [不知道是啥的数据包🤷] %v\n", err)
		return
	}
	log.Printf("[parseSCWebLiveWatchingUsers] [不知道是啥的数据包🤷] %s\n", jsonData)
}

func parseFeedPushPack(message []byte) {
	scWebFeedPush := &kuaishou_protostub.SCWebFeedPush{}
	if err := proto.Unmarshal(message, scWebFeedPush); err != nil {
		log.Printf("[parseFeedPushPack] [直播间弹幕🐎消息] %v\n", err)
		return
	}
	jsonData, err := json.Marshal(scWebFeedPush)
	if err != nil {
		log.Printf("[parseFeedPushPack] [直播间弹幕🐎消息] %v\n", err)
		return
	}
	log.Printf("[parseFeedPushPack] [直播间弹幕🐎消息] %s\n", jsonData)
}

func parseHeartBeatPack(message []byte) {
	heartAckMsg := &kuaishou_protostub.SCHeartbeatAck{}
	if err := proto.Unmarshal(message, heartAckMsg); err != nil {
		log.Printf("[parseHeartBeatPack] [心跳❤️响应] %v\n", err)
		return
	}
	jsonData, err := json.Marshal(heartAckMsg)
	if err != nil {
		log.Printf("[parseHeartBeatPack] [心跳❤️响应] %v\n", err)
		return
	}
	log.Printf("[parseHeartBeatPack] [心跳❤️响应] %s\n", jsonData)
}
