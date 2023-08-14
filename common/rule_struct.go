package common

import "live-room-crawler/constant"

type RuleResponse struct {
	DataList       []RuleGroupItem          `json:"dataList"`
	ResponseStatus *constant.ResponseStatus `json:"responseStatus,omitempty"`
	Count          int                      `json:"count,omitempty"`
}

type RuleGroupItem struct {
	RuleList                []Rule     `json:"ruleList"`
	Count                   int        `json:"count,omitempty"`
	ConditionTypeDictionary *Condition `json:"conditionTypeDictionary,omitempty"`
}

type Rule struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	UserID        string   `json:"userId,omitempty"`
	TenantID      string   `json:"tenantId,omitempty"`
	ProjectID     string   `json:"projectId,omitempty"`
	Sync          bool     `json:"sync,omitempty"`
	ReferID       int      `json:"referId,omitempty"`
	RuleType      int      `json:"ruleType,omitempty"`
	ConditionType string   `json:"conditionType,omitempty"`
	Threshold     int      `json:"threshold,omitempty"`
	DriverType    int      `json:"driverType,omitempty"`
	AnswerList    []Answer `json:"answerList,omitempty"`
}

type Condition struct {
	Code int    `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
}

type Answer struct {
	Text     string `json:"text"`
	AudioUrl string `json:"audioUrl"`
}
