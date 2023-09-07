package kuaishou

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *ConnectorStrategy) SendMsg(content string, liveStreamID string, color string) interface{} {
	variables := map[string]interface{}{
		"color":        color,
		"content":      content,
		"liveStreamId": liveStreamID,
	}
	query := `mutation SendLiveComment($liveStreamId: String, $content: String, $color: String) {
		sendLiveComment(liveStreamId: $liveStreamId, content: $content, color: $color) {
			result
			__typename
		}
	}`
	return c.liveGraphql("SendLiveComment", variables, query, map[string]string{})
}

func (c *ConnectorStrategy) Follow(principalID string, followType int) interface{} {
	variables := map[string]interface{}{
		"principalId": principalID,
		"type":        followType,
	}
	query := `mutation UserFollow($principalId: String, $type: Int) {
		webFollow(principalId: $principalId, type: $type) {
			followStatus
			__typename
		}
	}`
	return c.liveGraphql("UserFollow", variables, query, map[string]string{})
}

func (c *ConnectorStrategy) GetUserCardInfoByID(principalID string) interface{} {
	variables := map[string]interface{}{
		"principalId": principalID,
		"count":       3,
	}
	query := `query UserCardInfoById($principalId: String, $count: Int) {
		userCardInfo(principalId: $principalId, count: $count) {
			id
			originUserId
			avatar
			name
			description
			sex
			constellation
			cityName
			followStatus
			privacy
			feeds {
				eid
				photoId
				thumbnailUrl
				timestamp
				__typename
			}
			counts {
				fan
				follow
				photo
				__typename
			}
			__typename
		}
	}`
	return c.liveGraphql("UserCardInfoById", variables, query, map[string]string{})
}

func (c *ConnectorStrategy) GetAllGifts() interface{} {
	variables := map[string]interface{}{}
	query := `query AllGifts {
		allGifts
	}`
	return c.liveGraphql("AllGifts", variables, query, map[string]string{})
}

func (c *ConnectorStrategy) liveGraphql(operationName string, variables map[string]interface{}, query string, headers map[string]string) interface{} {

	data := map[string]interface{}{
		"operationName": operationName,
		"variables":     variables,
		"query":         query,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("[liveGraphql] [序列化失败]", err)
		return nil
	}

	// request room liveUrl
	req, err := http.NewRequest("POST", KuaishouApiHost, bytes.NewBuffer(payload))
	req.Header = c.Headers
	req.Header["Content-Type"] = []string{"application/json"}
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Create a new HTTP connection and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	var response interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Println("[liveGraphql] [解析失败]", err)
		return nil
	}

	log.Println("[liveGraphql] [操作返回数据] ｜", response)
	return response
}

func (c *ConnectorStrategy) unHex(data string) string {
	// 示例数据
	// data = 'E5 8C 97 E6 99 A8 E7 9A 84 E4 BF A1 E6 99 BA EF BC 88 E5 B7 B2 E7 B4 AB'
	data = strings.ReplaceAll(data, " ", "")
	decoded, err := hex.DecodeString(data)
	if err != nil {
		log.Println("[unHex] [解码失败]", err)
		return ""
	}
	return string(decoded)
}
