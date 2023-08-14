package data

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"live-room-crawler/common"
	"live-room-crawler/constant"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type RegistryItem struct {
	writeLock         sync.Mutex
	countLock         sync.Mutex
	conn              *websocket.Conn
	lostHeatBeatStamp int64
	StartRequest      common.CommandRequest
	RoomInfo          common.RoomInfo
	PlayDeque         *common.PlayDeque
	Statistics        map[common.CounterType]*common.StatisticCounter
	RuleGroupList     map[common.CounterType][]common.Rule
}

func (item *RegistryItem) CompareRule(counterType common.CounterType, counter *common.StatisticCounter) *common.CommandResponse {
	value, ok := item.Statistics[counterType]
	if ok {
		value.Add(counter)
	} else {
		return nil
	}

	lastMatch := value.LastMatch
	ruleList := item.RuleGroupList[counterType]

	for _, rule := range ruleList {
		if rule.AnswerList == nil || len(rule.AnswerList) == 0 {
			continue
		}

		if rule.Threshold > lastMatch && value.Count >= uint64(rule.Threshold) {
			idInt, _ := strconv.Atoi(rule.ID)
			value.LastMatch = rule.Threshold
			randomIndex := rand.Intn(len(rule.AnswerList))

			answer := rule.AnswerList[randomIndex]
			content := common.PlayContent{
				DrivenType: common.GetDrivenTypeByCode(rule.DriverType),
				Audio:      answer.AudioUrl,
				Text:       answer.Text,
			}

			response := common.CommandResponse{
				CommandType:    common.PLAY,
				TraceId:        uuid.NewString(),
				ResponseStatus: constant.SUCCESS,
				Content:        content,
				RuleMeta: common.RuleMeta{
					Id:        int64(idInt),
					Name:      rule.Name,
					Threshold: rule.Threshold,
					Type:      common.GUIDE,
				},
			}

			marshal, _ := json.Marshal(response)
			logger.Infof("☑️CompareRule match, push play message: %s", marshal)
			return &response
		}
	}

	return nil
}

func (item *RegistryItem) LoadRule() constant.ResponseStatus {
	item.countLock.Lock()
	defer item.countLock.Unlock()

	servicePart := item.StartRequest.Service
	responseStatus, ruleResponse := requestGetRule(servicePart)
	if !responseStatus.Success {
		return responseStatus
	}

	ruleRegistry := item.RuleGroupList
	for _, groupItem := range ruleResponse.DataList {
		dictionary := groupItem.ConditionTypeDictionary
		sort.Slice(groupItem.RuleList, func(i, j int) bool {
			return groupItem.RuleList[i].Threshold < groupItem.RuleList[j].Threshold
		})
		ruleRegistry[common.CounterType(dictionary.Name)] = groupItem.RuleList
	}

	return constant.SUCCESS
}

func (item *RegistryItem) UpdateStatistics(
	counterType common.CounterType,
	counter *common.StatisticCounter) {
	item.countLock.Lock()
	defer item.countLock.Unlock()
	marshal, _ := json.Marshal(counter)

	if item != nil {

		playResponse := item.CompareRule(counterType, counter)
		if playResponse != nil {
			item.WriteResponse(playResponse)
		}

		logger.Infof("UpdateStatistics for connection addr:%s type:%s counter:%s ", item.conn.RemoteAddr(), counterType, marshal)

	} else {
		logger.Warnf("UpdateStatistics FAIL for connection addr:%s", item.conn.RemoteAddr())
	}
}

func (item *RegistryItem) EnqueueAction(actionEvent common.UserActionEvent) {
	item.writeLock.Lock()
	defer item.writeLock.Unlock()
	marshal, _ := json.Marshal(actionEvent)
	logger.Infof("EnqueueAction invoked connection addr:%s event:%s", item.conn.RemoteAddr(), marshal)

	if item != nil {
		item.PlayDeque.PushBack(actionEvent)
	}
}

func (item *RegistryItem) DequeueAction() *common.UserActionEvent {
	item.writeLock.Lock()
	defer item.writeLock.Unlock()

	front := item.PlayDeque.PopFront()

	marshal, _ := json.Marshal(front)
	logger.Infof("DequeueAction invoked connection addr:%s event:%s", item.conn.RemoteAddr(), marshal)
	return front
}

func (item *RegistryItem) WriteResponse(response *common.CommandResponse) error {
	item.writeLock.Lock()
	defer item.writeLock.Unlock()
	marshal, err := json.Marshal(response)
	if err != nil {
		return err
	}
	if response.CommandType == common.PING {
		err := item.conn.WriteMessage(websocket.PongMessage, nil)
		return err
	} else {
		err := item.conn.WriteMessage(websocket.TextMessage, marshal)
		logger.Infof("🖋[EventDataRegistry]WriteResponse invoked connection addr:%s response:%s", item.conn.RemoteAddr(), marshal)
		return err
	}
}

func requestGetRule(servicePart common.ServiceStruct) (constant.ResponseStatus, *common.RuleResponse) {
	if servicePart.ApiBaseURL == "" {
		return constant.LOAD_RULE_FAIL, nil
	}

	params := url.Values{}
	params.Set("projectId", strconv.Itoa(servicePart.ProjectId))
	loadRuleUrl := strings.Join([]string{servicePart.ApiBaseURL, "/", constant.LoadGuideRuleURI, "?", params.Encode()}, "")

	req, err := http.NewRequest("GET", loadRuleUrl, nil)
	if err != nil {
		return constant.LOAD_RULE_FAIL, nil
	}
	// Set custom headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", servicePart.Authorization)
	req.Header.Set("TenantId", servicePart.TenantId)

	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	body := string(bodyBytes)

	if err != nil || response.StatusCode != 200 {
		return constant.LOAD_RULE_FAIL, nil
	}

	var ruleResponse common.RuleResponse
	logger.Infof("LoadRule result is: %s", body)
	err = json.Unmarshal(bodyBytes, &ruleResponse)
	return constant.SUCCESS, &ruleResponse
}
