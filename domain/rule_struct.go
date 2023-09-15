package domain

import (
	"live-room-crawler/constant"
)

type RuleResponse struct {
	DataList                        []RuleGroupItem                    `json:"dataList"`
	Project                         *Project                           `json:"project"`
	GloriousHonorPlatformCookieList []*GloriousHonorPlatformConfigItem `json:"gloriousHonorPlatformCookieList"`
	ResponseStatus                  *constant.ResponseStatus           `json:"responseStatus,omitempty"`
	Count                           int                                `json:"count,omitempty"`
}

type GloriousHonorPlatformConfigItem struct {
	Platform   string   `json:"platform"`
	Desc       string   `json:"desc"`
	CookieList []string `json:"cookieList"`
}

type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	UserID         string `json:"userId"`
	TenantID       string `json:"tenantId"`
	BackgroundID   string `json:"backgroundId"`
	DigitalHumanID string `json:"digitalHumanId"`
	NLPAppID       string `json:"nlpAppId"`
	NLPSkillID     string `json:"nlpSkillId"`
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
	Enable        bool     `json:"enable,omitempty"`
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

//func (configItem *GloriousHonorPlatformConfigItem) MarshalJSON() ([]byte, error) {
//	if configItem.CookieList == nil || len(configItem.CookieList) <= 0 {
//		return json.Marshal(configItem)
//	}
//
//	maskedCookies := make([]string, len(configItem.CookieList))
//	for i, cookie := range configItem.CookieList {
//		if len(cookie) > 10 {
//			maskedCookies[i] = fmt.Sprintf("%s***", cookie[:10])
//		} else {
//			maskedCookies[i] = cookie
//		}
//	}
//
//	type Alias GloriousHonorPlatformConfigItem
//	return json.Marshal(&struct {
//		*Alias
//		CookiesList []string `json:"cookieList"`
//	}{
//		Alias:       (*Alias)(configItem),
//		CookiesList: maskedCookies,
//	})
//}
