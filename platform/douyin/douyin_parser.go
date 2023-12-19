package douyin

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"io"
	"live-room-crawler/domain"
	"live-room-crawler/platform/douyin/douyin_protostub"
	"live-room-crawler/registry/data"
	"time"
)

func (c *ConnectorStrategy) OnMessage(message []byte, conn *websocket.Conn, localConn *websocket.Conn) {
	wssPackage := &douyin_protostub.PushFrame{}
	proto.Unmarshal(message, wssPackage)
	dataRegistry := data.GetDataRegistry()

	logId := wssPackage.LogId
	gzipReader, err := gzip.NewReader(bytes.NewReader(wssPackage.Payload))
	if err != nil {
		// Handle error
	}
	defer gzipReader.Close()

	payloadPackage := &douyin_protostub.Response{}
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
		select {
		case <-c.stopChan:
			// Stop signal received, exit the goroutine
			logger.Infof("StartListen was notified to stop by c.stopChan")
			return
		default:
			c.parseMessage(msg, dataRegistry, localConn)
		}

	}
}

func (c *ConnectorStrategy) parseMessage(msg *douyin_protostub.Message, dataRegistry *data.EventDataRegistry, localConn *websocket.Conn) {
	switch msg.Method {
	case "WebcastMatchAgainstScoreMessage":
		parseMatchAgainstScoreMessage(msg.Payload)
	case "WebcastLikeMessage":
		// 点赞消息WebcastLikeMessage；like,2
		// .total .count eg.	"total": 34136,
		likeMessage := parseWebcastLikeMessage(msg.Payload)
		dataRegistry.UpdateStatistics(localConn, domain.LIKE, domain.BuildStatisticsCounter(likeMessage.Total, false))
	case "WebcastMemberMessage":
		parseWebcastMemberMessage(msg.Payload)
		//dataRegistry.EnqueueAction(localConn, domain.UserActionEvent{
		//	Type:    domain.ON_ENTER,
		//	Username:  memberMessage.GetUser().NickName,
		//	Content:   memberMessage.GetUser().NickName + "加入了房间",
		//	EventTime: time.Unix(int64(memberMessage.Common.CreateTime), 0),
		//})
	case "WebcastGiftMessage":
		// 礼物消息WebcastGiftMessage; gift,3
		// .repeatCount 	"repeatCount": 10,
		// .comboCount	"comboCount": 10,
		// .domain.describe e.g 		"describe": "长孙明亮:送给主播 10个你最好看",
		giftMessage := parseWebcastGiftMessage(msg.Payload)
		dataRegistry.UpdateStatistics(localConn, domain.GIFT, domain.BuildStatisticsCounter(giftMessage.ComboCount, true))
		//dataRegistry.EnqueueAction(localConn, domain.UserActionEvent{
		//	Type:    domain.ON_GIFT,
		//	Username:  giftMessage.GetUser().NickName,
		//	Content:   giftMessage.Common.Describe,
		//	EventTime: time.Unix(int64(giftMessage.Common.CreateTime), 0),
		//})
	case "WebcastChatMessage":
		// comment,6
		chatMessage := parseWebcastChatMessage(msg.Payload)
		dataRegistry.UpdateStatistics(localConn, domain.COMMENT, domain.BuildStatisticsCounter(1, true))
		dataRegistry.EnqueueAction(localConn, domain.UserActionEvent{
			Type:      domain.COMMENT,
			Username:  chatMessage.GetUser().NickName,
			Content:   chatMessage.Content,
			EventTime: time.Unix(int64(chatMessage.EventTime), 0),
		})
	case "WebcastSocialMessage":
		// 关注WebcastSocialMessage； follow,4
		// .followCount, eg."followCount": "2193637",
		socialMessage := parseWebcastSocialMessage(msg.Payload)
		dataRegistry.UpdateStatistics(localConn, domain.FOLLOW, domain.BuildStatisticsCounter(socialMessage.FollowCount, true))
	case "WebcastRoomUserSeqMessage":
		//seqMessage := parseWebcastRoomUserSeqMessage(msg.Payload)
		// 榜单数据
		// .totalUser == totalPvForAnchor 浏览人数 view,5
		// .total == onlineUserForAnchor 在线人数 online,1
		seqMessage := parseWebcastRoomUserSeqMessage(msg.Payload)
		dataRegistry.UpdateStatistics(localConn, domain.VIEW, domain.BuildStatisticsCounter(uint64(seqMessage.TotalUser), true))
		dataRegistry.UpdateStatistics(localConn, domain.ONLINE, domain.BuildStatisticsCounter(uint64(seqMessage.Total), true))
	case "WebcastUpdateFanTicketMessage":
		parseWebcastUpdateFanTicketMessage(msg.Payload)
	case "WebcastCommonTextMessage":
		parseWebcastCommonTextMessage(msg.Payload)
	default:
		logger.Info("[onMessage] [⚠️" + msg.Method + "未知消息～]")
	}
}

func parseWebcastMemberMessage(payload []byte) *douyin_protostub.MemberMessage {
	chatMessage := &douyin_protostub.MemberMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseWebcastMemberMessage] [🏠加入房间消息] ｜ ", log)
	return chatMessage
}

func parseWebcastLikeMessage(payload []byte) *douyin_protostub.LikeMessage {
	chatMessage := &douyin_protostub.LikeMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseWebcastLikeMessage] [👍点赞消息] ｜ ", log)
	return chatMessage
}

func parseMatchAgainstScoreMessage(payload []byte) *douyin_protostub.MatchAgainstScoreMessage {
	chatMessage := &douyin_protostub.MatchAgainstScoreMessage{}
	proto.Unmarshal(payload, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseMatchAgainstScoreMessage] [📌MatchAgainstScoreMessage] ｜ ", log)
	return chatMessage
}

func parseWebcastChatMessage(data []byte) *douyin_protostub.ChatMessage {
	chatMessage := &douyin_protostub.ChatMessage{}
	proto.Unmarshal(data, chatMessage)
	jsonData, _ := json.Marshal(chatMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("⚛️⚛️⚛️✉️✉️✉️✉️✉️[parseWebcastChatMessage] [✉️直播间弹幕评论]｜", log)
	return chatMessage
}

func parseWebcastGiftMessage(data []byte) *douyin_protostub.GiftMessage {
	giftMessage := &douyin_protostub.GiftMessage{}
	proto.Unmarshal(data, giftMessage)
	jsonData, _ := json.Marshal(giftMessage)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonData, &dataMap)
	log := string(jsonData)
	logger.Info("[parseWebcastGiftMessage] [🎁直播间礼物] ｜ ", log)
	return giftMessage
}

func parseWebcastCommonTextMessage(data []byte) *douyin_protostub.CommonTextMessage {
	commonTextMessage := &douyin_protostub.CommonTextMessage{}
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

func parseWebcastUpdateFanTicketMessage(data []byte) *douyin_protostub.UpdateFanTicketMessage {
	updateFanTicketMessage := &douyin_protostub.UpdateFanTicketMessage{}
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

func parseWebcastRoomUserSeqMessage(data []byte) *douyin_protostub.RoomUserSeqMessage {
	roomUserSeqMessage := &douyin_protostub.RoomUserSeqMessage{}
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

func parseWebcastSocialMessage(data []byte) *douyin_protostub.SocialMessage {
	socialMessage := &douyin_protostub.SocialMessage{}
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
	obj := &douyin_protostub.PushFrame{
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
