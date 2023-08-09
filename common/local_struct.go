package common

import "live-room-crawler/constant"

type CommandType string

const (
	START CommandType = "start" // 开始直播；开始直播信号
	// LOAD // READY  CommandType = "ready" // 准备就绪； 直播间地址已配置就绪；
	LOAD CommandType = "load" // 更新规则
	STOP CommandType = "stop"
	PLAY CommandType = "play"
	PING CommandType = "ping" // ping响应
	PONG CommandType = "pong" // ping响应
)

type RuleType string

const (
	GUIDE RuleType = "guide"
	CHAT  RuleType = "chat"
)

type DrivenType string

const (
	TEXT  DrivenType = "text"
	AUDIO DrivenType = "audio"
)

type Platform string

const (
	DOUYIN   Platform = "douyin"   // 开始直播；开始直播信号
	KUAISHOU Platform = "kuaishou" // ping响应
)

// CommandRequest 请求体
type CommandRequest struct {
	CommandType CommandType   `json:"type"`
	Service     ServiceStruct `json:"service"`
	Target      TargetStruct  `json:"target"`
	TraceId     string        `json:"trace_id"`
}

type ServiceStruct struct {
	ApiBaseURL    string `json:"api_base_url"`
	ProjectId     string `json:"project_id"`
	Authorization string `json:"authorization"`
}

type TargetStruct struct {
	Platform Platform `json:"platform"`
	LiveURL  string   `json:"live_url"`
}

// CommandResponse 响应体
type CommandResponse struct {
	CommandType    CommandType             `json:"type,omitempty"`
	TraceId        string                  `json:"trace_id,omitempty"`
	Content        PlayContent             `json:"content,omitempty"`
	RuleMeta       RuleMeta                `json:"rule_meta,omitempty"`
	Room           RoomInfo                `json:"room,omitempty"`
	ResponseStatus constant.ResponseStatus `json:"status,omitempty"`
}

type PlayContent struct {
	DrivenType DrivenType `json:"driven_type,omitempty"`
	Text       string     `json:"text,omitempty"`
	Audio      string     `json:"audio,omitempty"`
}

type RuleMeta struct {
	Id        int64    `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Type      RuleType `json:"type,omitempty"`
	Threshold int32    `json:"threshold,omitempty"`
	Enable    bool     `json:"enable,omitempty"`
}

type RoomInfo struct {
	RoomId    string `json:"room_id,omitempty"`
	RoomTitle string `json:"title,omitempty"`
	Ttwid     string `json:"ttwid,omitempty"`
}
