package local_server

type CommandType string

const (
	START CommandType = "start" // 开始直播；开始直播信号
	// UPDATE // READY  CommandType = "ready" // 准备就绪； 直播间地址已配置就绪；
	UPDATE CommandType = "refresh" // 更新规则
	STOP   CommandType = "stop"    // 停止抓取
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
	ResponseType string                 `json:"type"`
	RoomTile     string                 `json:"room_tile"`
	TraceId      string                 `json:"trace_id"`
	Content      CrawlerResponseContent `json:"content"`
	RuleMeta     RuleMeta               `json:"rule_meta"`
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
}

type CrawlerResponseContent struct {
	TriggerType string `json:"trigger_type"`
	Text        string `json:"text"`
	Audio       string `json:"audio"`
	Enable      bool   `json:"enable"`
	TraceId     string `json:"trace_id"`
}

type RuleMeta struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Threshold int32  `json:"threshold"`
	Enable    bool   `json:"enable"`
}

type RoomInfo struct {
	RoomId    string
	RoomTitle string
	Ttwid     string
}
