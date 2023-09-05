package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/thoas/go-funk"
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
	stopChan          chan struct{} // Channel to signal stop
	writeLock         sync.Mutex
	countLock         sync.Mutex
	conn              *websocket.Conn
	lostHeatBeatStamp int64
	StartRequest      domain.CommandRequest
	RoomInfo          domain.RoomInfo
	PlayDeque         *util.PlayDeque
	Statistics        map[domain.CounterType]*domain.StatisticCounter
	RuleGroupList     map[domain.CounterType][]domain.Rule
	Project           *domain.Project
	ChatAvail         bool
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
		if !rule.Enable || rule.AnswerList == nil || len(rule.AnswerList) == 0 {
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
	if ruleResponse.Project != nil {
		item.Project = ruleResponse.Project
		item.ChatAvail = ruleResponse.Project.NLPAppID != ""
		logger.Infof("⚙️LoadRule update ChatAvail: %b, NLPAppId: %s", item.ChatAvail, ruleResponse.Project.NLPAppID)
	}

	ruleRegistry := item.RuleGroupList
	for _, groupItem := range ruleResponse.DataList {
		dictionary := groupItem.ConditionTypeDictionary
		// filter out all enabled rules
		filteredList := funk.Filter(groupItem.RuleList, func(rule domain.Rule) bool {
			return rule.Enable
		}).([]domain.Rule)
		// sort by Threshold, DriverType desc
		sort.Slice(filteredList, func(i, j int) bool {
			return filteredList[i].Threshold >= filteredList[j].Threshold &&
				filteredList[i].DriverType > filteredList[j].DriverType
		})

		// load rules which only enabled
		ruleRegistry[domain.CounterType(dictionary.Name)] = filteredList
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

func (item *RegistryItem) RequestNlp(username string, content string) *domain.QueryResponse {
	item.writeLock.Lock()
	defer item.writeLock.Unlock()
	logger.Infof("🖥☎️[data.RegistryItem] RequestNlp enter for user:%s query: %s", username, content)

	servicePart := item.StartRequest.Service
	startId := item.StartRequest.TraceId
	requestURL := strings.Join([]string{servicePart.ApiBaseURL, "/", constant.QueryNlpURI}, "")

	projectId, _ := strconv.Atoi(servicePart.ProjectId)
	requestBody := domain.QueryRequest{
		ProjectID: projectId,
		Query:     content,
		SessionID: strings.Join([]string{startId, username}, "-"),
		TraceID:   strings.Join([]string{startId, uuid.NewString()}, "-"),
	}

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	logger.Infof("🖥📞[data.RegistryItem] RequestNlp request traceId:%s with:%s", requestBody.TraceID, jsonBody)

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", servicePart.Authorization)
	req.Header.Set("TenantId", servicePart.TenantId)
	req.Header.Set("CustomTraceId", startId+uuid.NewString())

	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	bodyBytes, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	responseBody := bodyBytes

	logger.Infof("🖥📡[data.RegistryItem] RequestNlp response traceId:%s code:%d with:%s", requestBody.TraceID, response.StatusCode, responseBody)
	queryResponse := &domain.QueryResponse{}
	json.Unmarshal(responseBody, queryResponse)
	if err != nil || !util.IsInList(response.StatusCode, []int{200, 401, 403, 400}) {
		queryResponse.ResponseStatus = constant.UNKNOWN_NLP_RESPONSE
		return queryResponse
	}

	queryResponse.Meta.UserName = username
	return queryResponse
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
		return constant.REQUEST_SERVER_FAILED, nil
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
	if err != nil || !util.IsInList(response.StatusCode, []int{200, 401, 403, 400}) {
		return constant.REQUEST_SERVER_FAILED, nil
	}

	var ruleResponse domain.RuleResponse
	logger.Infof("LoadRule result is: %s", body)
	err = json.Unmarshal(bodyBytes, &ruleResponse)
	if err != nil || ruleResponse.ResponseStatus == nil {
		return constant.LOAD_RULE_FAIL, nil
	}

	return *ruleResponse.ResponseStatus, &ruleResponse
}
