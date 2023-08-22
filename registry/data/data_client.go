package data

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"live-room-crawler/constant"
	"live-room-crawler/domain"
	"live-room-crawler/util"
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
	StartRequest      domain.CommandRequest
	RoomInfo          domain.RoomInfo
	PlayDeque         *util.PlayDeque
	Statistics        map[domain.CounterType]*domain.StatisticCounter
	RuleGroupList     map[domain.CounterType][]domain.Rule
}

func (item *RegistryItem) CompareRule(counterType domain.CounterType, counter *domain.StatisticCounter) *domain.CommandResponse {
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
			content := domain.PlayContent{
				DrivenType: domain.GetDrivenTypeByCode(rule.DriverType),
				Audio:      answer.AudioUrl,
				Text:       answer.Text,
			}

			response := domain.CommandResponse{
				CommandType:    domain.PLAY,
				TraceId:        uuid.NewString(),
				ResponseStatus: constant.SUCCESS,
				Content:        content,
				RuleMeta: domain.RuleMeta{
					Id:        int64(idInt),
					Name:      rule.Name,
					Threshold: rule.Threshold,
					Type:      domain.GUIDE,
				},
			}

			marshal, _ := json.Marshal(response)
			logger.Infof("☑️CompareRule match, push play message: %s", marshal)
			return &response
		}
	}

	return nil
}

func (item *RegistryItem) LoadRule(traceId string) constant.ResponseStatus {
	item.countLock.Lock()
	defer item.countLock.Unlock()

	servicePart := item.StartRequest.Service
	responseStatus, ruleResponse := requestGetRule(traceId, servicePart)
	if !responseStatus.Success {
		return responseStatus
	}

	if ruleResponse.DataList != nil && len(ruleResponse.DataList) > 0 { // ? how if user clear his rules
		item.RuleGroupList = make(map[domain.CounterType][]domain.Rule)
	}
	ruleRegistry := item.RuleGroupList
	for _, groupItem := range ruleResponse.DataList {
		dictionary := groupItem.ConditionTypeDictionary
		sort.Slice(groupItem.RuleList, func(i, j int) bool {
			return groupItem.RuleList[i].Threshold > groupItem.RuleList[j].Threshold &&
				groupItem.RuleList[i].DriverType > groupItem.RuleList[j].DriverType
		})
		ruleRegistry[domain.CounterType(dictionary.Name)] = groupItem.RuleList
	}

	return constant.SUCCESS
}

func (item *RegistryItem) UpdateStatistics(
	counterType domain.CounterType,
	counter *domain.StatisticCounter) {
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

func (item *RegistryItem) EnqueueAction(actionEvent domain.UserActionEvent) {
	item.writeLock.Lock()
	defer item.writeLock.Unlock()
	marshal, _ := json.Marshal(actionEvent)
	logger.Infof("EnqueueAction invoked connection addr:%s event:%s", item.conn.RemoteAddr(), marshal)

	if item != nil {
		item.PlayDeque.PushBack(actionEvent)
	}
}

func (item *RegistryItem) DequeueAction() *domain.UserActionEvent {
	item.writeLock.Lock()
	defer item.writeLock.Unlock()

	front := item.PlayDeque.PopFront()

	marshal, _ := json.Marshal(front)
	logger.Infof("DequeueAction invoked connection addr:%s event:%s", item.conn.RemoteAddr(), marshal)
	return front
}

func (item *RegistryItem) WriteResponse(response *domain.CommandResponse) error {
	item.writeLock.Lock()
	defer item.writeLock.Unlock()
	marshal, err := json.Marshal(response)
	if err != nil {
		return err
	}
	if response.CommandType == domain.PING {
		err := item.conn.WriteMessage(websocket.PongMessage, nil)
		return err
	} else {
		err := item.conn.WriteMessage(websocket.TextMessage, marshal)
		logger.Infof("🖋[EventDataRegistry]WriteResponse invoked connection addr:%s response:%s", item.conn.RemoteAddr(), marshal)
		return err
	}
}

func requestGetRule(traceId string, servicePart domain.ServiceStruct) (constant.ResponseStatus, *domain.RuleResponse) {
	if servicePart.ApiBaseURL == "" {
		return constant.LOAD_RULE_FAIL, nil
	}

	params := url.Values{}
	params.Set("projectId", servicePart.ProjectId)
	loadRuleUrl := strings.Join([]string{servicePart.ApiBaseURL, "/", constant.LoadGuideRuleURI, "?", params.Encode()}, "")

	req, err := http.NewRequest("GET", loadRuleUrl, nil)
	if err != nil {
		return constant.LOAD_RULE_FAIL, nil
	}
	// Set custom headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", servicePart.Authorization)
	req.Header.Set("TenantId", servicePart.TenantId)
	req.Header.Set("CustomTraceId", traceId)

	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	body := string(bodyBytes)

	logger.Infof("LoadRule request url:%s with code %d result is: %s", loadRuleUrl, response.StatusCode, body)
	if err != nil || response.StatusCode != 200 {
		return constant.LOAD_RULE_FAIL, nil
	}

	var ruleResponse domain.RuleResponse
	logger.Infof("LoadRule result is: %s", body)
	err = json.Unmarshal(bodyBytes, &ruleResponse)
	if err != nil || ruleResponse.ResponseStatus == nil {
		return constant.LOAD_RULE_FAIL, nil
	}

	return constant.SUCCESS, &ruleResponse
}
