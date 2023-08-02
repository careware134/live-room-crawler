package local_server

type CommandType string

const (
	START CommandType = "start" // 开始直播；开始直播信号
	// LOAD // READY  CommandType = "ready" // 准备就绪； 直播间地址已配置就绪；
	LOAD CommandType = "refresh" // 更新规则
	STOP CommandType = "stop"
	PLAY CommandType = "play"
	PONG CommandType = "pong" // ping响应
)

type CommandRequest struct {
	CommandType   CommandType `json:"type"`
	BaseURL       string      `json:"api_base_url"`
	ProjectId     string      `json:"project_id"`
	Authorization string      `json:"authorization"`
	LiveURL       string      `json:"live_url"`
	TraceId       string      `json:"trace_id"`
}

type ResponseType string

const (
	GUIDE ResponseType = "guide"
	CHAT  ResponseType = "chat"
)

type ContentType string

const (
	TEXT  ContentType = "text"
	AUDIO ContentType = "audio"
)

type CommandResponse struct {
	CommandType CommandType            `json:"type,omitempty"`
	DrivenType  string                 `json:"driven_type,omitempty"`
	RoomTile    string                 `json:"room_tile,omitempty"`
	TraceId     string                 `json:"trace_id,omitempty"`
	Content     CrawlerResponseContent `json:"content,omitempty"`
	RuleMeta    RuleMeta               `json:"rule_meta,omitempty"`
	Success     bool                   `json:"success,omitempty"`
	Message     string                 `json:"message,omitempty"`
}

type CrawlerResponseContent struct {
	TriggerType ResponseType `json:"trigger_type,omitempty"`
	Text        string       `json:"text,omitempty"`
	Audio       string       `json:"audio,omitempty"`
}

type RuleMeta struct {
	Id        int64  `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	Threshold int32  `json:"threshold,omitempty"`
	Enable    bool   `json:"enable,omitempty"`
}

type RoomInfo struct {
	RoomId    string
	RoomTitle string
	Ttwid     string
}
