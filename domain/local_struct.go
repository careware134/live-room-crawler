package domain

import (
	"live-room-crawler/constant"
)

type CommandType string

const (
	START CommandType = "start" // 开始直播；开始直播信号
	// READY  CommandType = "ready" // 准备就绪； 直播间地址已配置就绪；
	// LOAD    CommandType = "load"    // 更新规则
	REFRESH CommandType = "refresh" // 更新规则
	STOP    CommandType = "stop"
	EVENT   CommandType = "event" // headless模式，PC_LIVE通过headless-chrome抓取数据反推评论数据
	PLAY    CommandType = "play"
	PING    CommandType = "ping" // ping响应
	PONG    CommandType = "pong" // ping响应
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

var drivenTypeMap = map[int]DrivenType{
	1: TEXT,
	2: AUDIO,
}

func GetDrivenTypeByCode(index int) DrivenType {
	drivenType, ok := drivenTypeMap[index]
	if ok {
		return drivenType
	}

	return TEXT
}

type Platform string

const (
	DOUYIN    Platform = "douyin"
	KUAISHOU  Platform = "kuaishou"
	MEITUAN   Platform = "meituan"   // 开始直播；开始直播信号
	PINDUODUO Platform = "pinduoduo" // 开始直播；开始直播信号
)

type PlayMode string

const (
	HOST_MODE   PlayMode = "HOST_MODE"
	ASSIST_MODE PlayMode = "ASSIST_MODE"
)

var playModeMap = map[int]PlayMode{
	1: HOST_MODE,
	2: ASSIST_MODE,
}

func GetPlayModeByCode(index int) PlayMode {
	playMode, ok := playModeMap[index]
	if ok {
		return playMode
	}

	return HOST_MODE
}

// CommandRequest 请求体
type CommandRequest struct {
	CommandType CommandType      `json:"type"`
	Service     *ServiceStruct   `json:"service,omitempty"`
	Target      *TargetStruct    `json:"target,omitempty"`
	RoomInfo    *RoomInfo        `json:"room,omitempty"`
	TraceId     string           `json:"trace_id"`
	ActionEvent *UserActionEvent `json:"event,omitempty"`
}

type ServiceStruct struct {
	ApiBaseURL    string `json:"api_base_url"`
	ProjectId     string `json:"project_id"`
	TenantId      string `json:"tenant_id"`
	Authorization string `json:"authorization"`
}

type TargetStruct struct {
	Platform Platform `json:"platform"`
	LiveURL  string   `json:"live_url"`
	Cookie   string   `json:"cookie"`
	Headless bool     `json:"headless"`
}

// CommandResponse 响应体
type CommandResponse struct {
	CommandType    CommandType             `json:"type,omitempty"`
	TraceId        string                  `json:"trace_id,omitempty"`
	Content        *PlayContent            `json:"content,omitempty"`
	RuleMeta       *RuleMeta               `json:"rule_meta,omitempty"`
	Room           *RoomInfo               `json:"room,omitempty"`
	ResponseStatus constant.ResponseStatus `json:"status,omitempty"`
}

type PlayContent struct {
	DrivenType DrivenType `json:"driven_type,omitempty"`
	PlayMode   PlayMode   `json:"play_mode,omitempty"`
	Text       string     `json:"text,omitempty"`
	Audio      string     `json:"audio,omitempty"`
}

type RuleMeta struct {
	Id        int64    `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Query     string   `json:"query,omitempty"`
	UserName  string   `json:"user_name,omitempty"`
	Type      RuleType `json:"type,omitempty"`
	Threshold int      `json:"threshold,omitempty"`
	Enable    bool     `json:"enable,omitempty"`
}

type RoomInfo struct {
	RoomId       string `json:"room_id,omitempty"`
	RoomTitle    string `json:"title,omitempty"`
	Token        string `json:"token,omitempty"`
	WebSocketUrl string `json:"web_socket_url,omitempty"`
}
