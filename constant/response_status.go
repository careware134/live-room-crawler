package constant

type ResponseStatus struct {
	Success bool   `json:"success,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

var (
	SUCCESS          = ResponseStatus{Success: true, Code: "SUCCESS", Message: "成功"}
	CONNECTION_FAIL  = ResponseStatus{Success: false, Code: "CONNECTION_FAIL", Message: "直播平台连接失败！"}
	INVALID_LIVE_URL = ResponseStatus{Success: false, Code: "INVALID_LIVE_URL", Message: "无效直播间地址，请确认！"}
)
