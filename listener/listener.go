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
			unPackMatchAgainstScoreMessage(msg.Payload)
		case "WebcastLikeMessage":
			unPackWebcastLikeMessage(msg.Payload)
		case "WebcastMemberMessage":
			unPackWebcastMemberMessage(msg.Payload)
		case "WebcastGiftMessage":
			unPackWebcastGiftMessage(msg.Payload)
		case "WebcastChatMessage":
			unPackWebcastChatMessage(msg.Payload)
		case "WebcastSocialMessage":
			unPackWebcastSocialMessage(msg.Payload)
		case "WebcastRoomUserSeqMessage":
			unPackWebcastRoomUserSeqMessage(msg.Payload)
		case "WebcastUpdateFanTicketMessage":
			unPackWebcastUpdateFanTicketMessage(msg.Payload)
		case "WebcastCommonTextMessage":
			unPackWebcastCommonTextMessage(msg.Payload)
		default:
			logger.Info("[onMessage] [⚠️" + msg.Method + "未知消息～]")
		}
	}
}

func unPackWebcastMemberMessage(payload []byte) {
	chatMessage := &protostub.MemberMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Infof("[unPackWebcastMemberMessage] [🚹🚺直播间成员加入消息] ｜ " + log)
}

func unPackWebcastLikeMessage(payload []byte) {
	chatMessage := &protostub.LikeMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[unPackWebcastLikeMessage] [👍直播间点赞消息]" + log)
}

func unPackMatchAgainstScoreMessage(payload []byte) {
	chatMessage := &protostub.MatchAgainstScoreMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[unPackMatchAgainstScoreMessage] [🤷不知道是啥的消息] ｜ " + log)
}

func unPackWebcastChatMessage(data []byte) map[string]interface{} {
	chatMessage := &protostub.ChatMessage{}
	proto.Unmarshal(data, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[unPackWebcastChatMessage] [📧直播间弹幕消息]｜ %s", log)
	return dataMap
}

func unPackWebcastGiftMessage(data []byte) {
	giftMessage := &protostub.GiftMessage{}
	proto.Unmarshal(data, giftMessage)
	jsonData, _ := json.Marshal(giftMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[unPackWebcastGiftMessage] [🎁直播间礼物消息] ｜ " + log)
}

func unPackWebcastCommonTextMessage(data []byte) {
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

	logger.Infof("[unPackWebcastCommonTextMessage] | %s", jsonStr)
}

func unPackWebcastUpdateFanTicketMessage(data []byte) map[string]interface{} {
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

	logger.Info("[unPackWebcastUpdateFanTicketMessage]｜ " + jsonStr)
	return dataMap
}

func unPackWebcastRoomUserSeqMessage(data []byte) map[string]interface{} {
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

	var dataMap map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &dataMap)
	if err != nil {
		// Handle error
	}

	logger.Infof("[unPackWebcastRoomUserSeqMessage] [️🏄🏂用户信息]｜ " + jsonStr)
	return dataMap
}

func unPackWebcastSocialMessage(data []byte) map[string]interface{} {
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

	logger.Infof("[unPackWebcastSocialMessage] [➕直播间关注消息] ｜ " + jsonStr)
	return dataMap
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
