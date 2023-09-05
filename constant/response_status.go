package constant

type ResponseStatus struct {
	Success bool   `json:"success"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

var (
	SUCCESS                = ResponseStatus{Success: true, Code: "SUCCESS", Message: "成功"}
	UNKNOWN_COMMAND        = ResponseStatus{Success: false, Code: "UNKNOWN_COMMAND", Message: "未知指令！"}
	SUCCESS_ALREADY        = ResponseStatus{Success: true, Code: "SUCCESS_ALREADY", Message: "成功，请勿重复start！"}
	LOAD_RULE_FAIL         = ResponseStatus{Success: false, Code: "LOAD_RULE_FAIL", Message: "加载规则失败！"}
	REQUEST_SERVER_FAILED  = ResponseStatus{Success: false, Code: "REQUEST_SERVER_FAILED", Message: "网络请求失败！"}
	CLIENT_NOT_READY       = ResponseStatus{Success: false, Code: "CLIENT_NOT_READY", Message: "客户端未准备就绪，请检查配置！"}
	CONNECTION_FAIL        = ResponseStatus{Success: false, Code: "CONNECTION_FAIL", Message: "直播平台连接失败！"}
	LIVE_CONNECTION_CLOSED = ResponseStatus{Success: false, Code: "LIVE_CONNECTION_CLOSED", Message: "连接回收：心跳丢失或直播间关闭！"}
	INVALID_LIVE_URL       = ResponseStatus{Success: false, Code: "INVALID_LIVE_URL", Message: "无效直播间地址，请确认！"}
	UNKNOWN_PLATFORM       = ResponseStatus{Success: false, Code: "UNKNOWN_PLATFORM", Message: "未知直播间平台！"}

	UNKNOWN_NLP_RESPONSE = ResponseStatus{Success: false, Code: "UNKNOWN_NLP_RESPONSE", Message: "请求NLP返回未知响应！"}
)

const (
	PlayUserAction          = false
	PlayDequeuePushInterval = 1
	HeartbeatCheckInterval  = 150
	LogRound                = 60
	LoadGuideRuleURI        = "rule/guide/loadByProjectId"
	QueryNlpURI             = "rule/chat/model/reply"
)
