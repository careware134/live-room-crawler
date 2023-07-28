package listener

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"io"
	"live-room-crawler/protostub"
	"live-room-crawler/util"
)

var logger = util.Logger()

func OnMessage(message []byte, conn *websocket.Conn) {
	wssPackage := &protostub.PushFrame{}
	proto.Unmarshal(message, wssPackage)

	logId := wssPackage.LogId
	gzipReader, err := gzip.NewReader(bytes.NewReader(wssPackage.Payload))
	if err != nil {
		// Handle error
	}
	defer gzipReader.Close()

	payloadPackage := &protostub.Response{}
	data, err := io.ReadAll(gzipReader)

	err = proto.Unmarshal(data, payloadPackage)
	if err != nil {
		// Handle error
	}

	// 发送ack包
	if payloadPackage.NeedAck {
		sendAck(conn, logId, payloadPackage.InternalExt)
	}

	for _, msg := range payloadPackage.MessagesList {
		switch msg.Method {
		case "WebcastMatchAgainstScoreMessage":
			parseMatchAgainstScoreMessage(msg.Payload)
		case "WebcastLikeMessage":
			parseWebcastLikeMessage(msg.Payload)
		case "WebcastMemberMessage":
			parseWebcastMemberMessage(msg.Payload)
		case "WebcastGiftMessage":
			parseWebcastGiftMessage(msg.Payload)
		case "WebcastChatMessage":
			parseWebcastChatMessage(msg.Payload)
		case "WebcastSocialMessage":
			parseWebcastSocialMessage(msg.Payload)
		case "WebcastRoomUserSeqMessage":
			//seqMessage := parseWebcastRoomUserSeqMessage(msg.Payload)
			parseWebcastRoomUserSeqMessage(msg.Payload)
			// TODO update local registry
		case "WebcastUpdateFanTicketMessage":
			parseWebcastUpdateFanTicketMessage(msg.Payload)
		case "WebcastCommonTextMessage":
			parseWebcastCommonTextMessage(msg.Payload)
		default:
			logger.Info("[onMessage] [⚠️" + msg.Method + "未知消息～]")
		}
	}
}

func parseWebcastMemberMessage(payload []byte) {
	chatMessage := &protostub.MemberMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseWebcastMemberMessage] [🏠加入房间消息] ｜ ", log)
}

func parseWebcastLikeMessage(payload []byte) *protostub.LikeMessage {
	chatMessage := &protostub.LikeMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseWebcastLikeMessage] [👍点赞消息] ｜ ", log)
	return chatMessage
}

func parseMatchAgainstScoreMessage(payload []byte) *protostub.MatchAgainstScoreMessage {
	chatMessage := &protostub.MatchAgainstScoreMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseMatchAgainstScoreMessage] [📌MatchAgainstScoreMessage] ｜ ", log)
	return chatMessage
}

func parseWebcastChatMessage(data []byte) *protostub.ChatMessage {
	chatMessage := &protostub.ChatMessage{}
	proto.Unmarshal(data, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseWebcastChatMessage] [✉️直播间弹幕评论]｜", log)
	return chatMessage
}

func parseWebcastGiftMessage(data []byte) *protostub.GiftMessage {
	giftMessage := &protostub.GiftMessage{}
	proto.Unmarshal(data, giftMessage)
	jsonData, _ := json.Marshal(giftMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseWebcastGiftMessage] [🎁直播间礼物] ｜ ", log)
	return giftMessage
}

func parseWebcastCommonTextMessage(data []byte) *protostub.CommonTextMessage {
	commonTextMessage := &protostub.CommonTextMessage{}
	err := proto.Unmarshal(data, commonTextMessage)
	if err != nil {
		// Handle error
	}

	marshaler := &jsonpb.Marshaler{OrigName: true}
	jsonStr, err := marshaler.MarshalToString(commonTextMessage)
	if err != nil {
		// Handle error
	}

	logger.Info("[parseWebcastCommonTextMessage] |", jsonStr)
	return commonTextMessage
}

func parseWebcastUpdateFanTicketMessage(data []byte) *protostub.UpdateFanTicketMessage {
	updateFanTicketMessage := &protostub.UpdateFanTicketMessage{}
	err := proto.Unmarshal(data, updateFanTicketMessage)
	if err != nil {
		// Handle error
	}

	marshaler := &jsonpb.Marshaler{OrigName: true}
	jsonStr, err := marshaler.MarshalToString(updateFanTicketMessage)
	if err != nil {
		// Handle error
	}

	var dataMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &dataMap)
	if err != nil {
		// Handle error
	}

	logger.Info("[parseWebcastUpdateFanTicketMessage] [💝粉丝数更新消息]｜ ", jsonStr)
	return updateFanTicketMessage
}

func parseWebcastRoomUserSeqMessage(data []byte) *protostub.RoomUserSeqMessage {
	roomUserSeqMessage := &protostub.RoomUserSeqMessage{}
	err := proto.Unmarshal(data, roomUserSeqMessage)
	if err != nil {
		// Handle error
	}

	marshaler := &jsonpb.Marshaler{OrigName: true}
	jsonStr, err := marshaler.MarshalToString(roomUserSeqMessage)
	if err != nil {
		// Handle error
	}

	logger.Info("[parseWebcastRoomUserSeqMessage] [️🏂用户榜单信息]｜ ", jsonStr)
	return roomUserSeqMessage
}

func parseWebcastSocialMessage(data []byte) *protostub.SocialMessage {
	socialMessage := &protostub.SocialMessage{}
	err := proto.Unmarshal(data, socialMessage)
	if err != nil {
		// Handle error
	}

	marshaler := &jsonpb.Marshaler{OrigName: true}
	jsonStr, err := marshaler.MarshalToString(socialMessage)
	if err != nil {
		// Handle error
	}

	var dataMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &dataMap)
	if err != nil {
		// Handle error
	}

	logger.Info("[parseWebcastSocialMessage] [➕直播间关注消息] ｜ ", jsonStr)
	return socialMessage
}

func sendAck(ws *websocket.Conn, logId uint64, internalExt string) {
	obj := &protostub.PushFrame{
		//PayloadType: "ack",
		LogId:       logId,
		PayloadType: internalExt,
	}

	data, err := proto.Marshal(obj)
	if err != nil {
		// Handle error
	}

	err = ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		// Handle error
	}
}
